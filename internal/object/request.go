package object

import (
	"fmt"
	"reflect"

	"github.com/seerx/chain/internal/object/validators"

	"github.com/seerx/chain/internal/reflects"
)

type RequestObjectField struct {
	Name         string // json tag 或者 fieldName，用于从  map 中获取值
	FieldName    string // 结构字段名称
	Type         reflect.Type
	Ptr          bool                   // 定义的类型是否是指针
	Slice        bool                   // 定义的类型是都是切片
	SliceType    reflect.Type           // 切片类型定义
	SliceItemPtr bool                   // Slice 项的类型是否是指针
	Require      bool                   // 必填参数
	Validators   []validators.Validator // 数据验证
}

func GenerateRequestObjectField(name string, fieldName string, info *reflects.TypeInfo, require bool) *RequestObjectField {
	field := &RequestObjectField{
		Name:      name,
		FieldName: fieldName,
		Type:      info.Reference,
		Ptr:       info.IsRawPtr,
		Slice:     info.IsRawSlice,
		SliceType: nil,
		//SliceType:    info.,
		SliceItemPtr: info.IsSliceItemIsPtr,
		Require:      require,
		Validators:   nil,
	}

	if field.Slice {
		field.SliceType = info.Raw
	}

	return field
}

type RequestObject struct {
	TypeName string       // 类型名称
	Type     reflect.Type // 实际数据类型

	Primitive bool // 原生类型
	Struct    bool //结构

	Fields []*RequestObjectField // struct.field 结构中定义的字段
	//Children   []*InputVar            // 结构属性类型
	//Validators []validators.Validator // 数据有效性检查
}

func (rof *RequestObjectField) NewInstance(parentPath string, data interface{}, mgr *RequestObjectManager) (reflect.Value, error) {
	if data == nil && rof.Require {
		// 不能为空
		return reflect.ValueOf(nil), fmt.Errorf("%s.%s is required", parentPath, rof.Name)
	}
	objType := mgr.Find(rof.Type)
	if objType == nil {
		return reflect.ValueOf(nil), fmt.Errorf("No RequestObject exists: %s", rof.Name)
	}

	if rof.Slice {
		// 切片
		if data == nil {
			// 数据是 nil ，返回空数组
			return reflect.MakeSlice(rof.SliceType, 0, 0), nil
		}
		ary, ok := data.([]interface{})
		if !ok {
			return reflect.ValueOf(nil), fmt.Errorf("Cann't parse %s as slice", rof.Name)
		}

		itemObj := mgr.Find(rof.Type)
		if itemObj == nil {
			return reflect.ValueOf(nil), fmt.Errorf("Cann't find %s's object'", rof.Name)
		}
		slice := reflect.MakeSlice(rof.SliceType, 0, len(ary))
		for _, v := range ary {
			if item, err := itemObj.NewInstance(parentPath, rof.Name, v, mgr); err != nil {
				return reflect.ValueOf(nil), err
			} else {
				if rof.SliceItemPtr {
					// 元素是指针
					slice = reflect.Append(slice, item)
				} else {
					// 元素非指针
					slice = reflect.Append(slice, item.Elem())
				}
			}
		}

		// 返回切片
		return slice, nil
	}
	// 非切片类型
	val, err := objType.NewInstance(parentPath, rof.Name, data, mgr)
	if err != nil {
		return reflect.ValueOf(nil), err
	}

	// 验证数据合法性
	for _, vld := range rof.Validators {
		if err := vld.Check(val.Interface()); err != nil {
			// 数据校验不通过
			return reflect.ValueOf(nil), err
		}
	}

	if rof.Ptr {
		// 指针
		return val, nil
	} else {
		// 非指针
		return val.Elem(), nil
	}
	//return nil, nil
}

func (ro *RequestObject) NewInstance(parentPath string, fieldName string, data interface{}, mgr *RequestObjectManager) (reflect.Value, error) {
	if data == nil {
		// 数据是空的
		return reflect.ValueOf(nil), nil
	}
	if ro.Primitive {
		// 原生类型
		// 对 data 做类型判断及数据转换
		outData := tryToConvert(ro.Type, data)
		return outData, nil
	}

	//if ro.Struct {
	// 结构类型
	//}

	mp, ok := data.(map[string]interface{})
	if !ok {
		panic(fmt.Errorf("Cann't parse %s as struct", ro.TypeName))
	}
	inst := reflect.New(ro.Type)
	elem := inst.Elem()
	var thisParent string
	if parentPath == "" {
		thisParent = fieldName
	} else {
		thisParent = parentPath + "." + fieldName
	}

	for _, fd := range ro.Fields {
		v, ok := mp[fd.Name]

		if ok {
			field := elem.FieldByName(fd.FieldName)
			//obj := mgr.Find(fd.Type)
			objVal, err := fd.NewInstance(thisParent, v, mgr)
			//objVal, err := obj.NewInstance(thisParent, fd.Name, v, mgr)
			if err != nil {
				return reflect.ValueOf(nil), err
			}
			if fd.Ptr {
				field.Set(objVal)
			} else {
				field.Set(objVal) // .Elem())
			}
			// TODO 添加到已发现字段，用于 reqiure 函数判断
		} else {
			if fd.Require {
				// 必填字段
				return reflect.ValueOf(nil), fmt.Errorf("%s.%s is required", parentPath, fd.Name)
			}
		}
	}

	return inst, nil
}

package schema

import (
	"fmt"
	"reflect"

	"github.com/seerx/chain/internal/reflects"

	"github.com/seerx/chain/internal/object"
	"github.com/seerx/chain/pkg/apimap"
)

// traversal 遍历 typ
// 如果是遍历请求参数，则 inputMap 必须不为 nil ，并且会返回得到的 InputObject
// 如果是便器范库参数，则 inputMap 必须为 nil
// 如果 返回的 error 不为 nil，则说明参数有不合法的内容
func traversal(loc *reflects.Location,
	typ reflect.Type,
	referenceMap map[string]int,
	objMap map[string]*apimap.ObjectInfo,
	requestMgr *object.RequestObjectManager) (*apimap.ObjectInfo, *object.RequestObject, error) {

	tp := reflects.ParseType(loc, typ)
	// 指向类型是结构体
	obj, exists := objMap[tp.ID()]
	if exists {
		obj.ReferenceCount++
		if rc, e := referenceMap[obj.ID]; e {
			referenceMap[obj.ID] = rc + 1
		} else {
			referenceMap[obj.ID] = 1
		}
		// 返回指向结构体
		refObj := &apimap.ObjectInfo{
			ID:   obj.ID,
			Type: obj.Type,
		}

		if requestMgr != nil {
			// 从输入对象中查找，并一块返回
			ro := requestMgr.Find(tp.Reference)
			if ro == nil {
				return refObj, nil, fmt.Errorf("No Request Object found [%s]: %s", tp.Name, loc)
			}
			return refObj, ro, nil
		}
		return refObj, nil, nil
	}
	// 还没有在 info.Objects 中注册，顶级对象不要 name require 等属性，只需要 ID，Type 和 Children
	obj = &apimap.ObjectInfo{
		ID:       tp.ID(),
		Type:     tp.TypeName(),
		Children: nil,
	}
	objMap[obj.ID] = obj
	var requestObj *object.RequestObject
	if requestMgr != nil {
		// 生成 RequestObject 顶级对象
		requestObj = &object.RequestObject{
			TypeName:  tp.Name,
			Type:      tp.Reference,
			Primitive: tp.IsPrimitive,
			Struct:    tp.IsStruct,
			Fields:    nil,
		}
		requestMgr.Register(requestObj)
	}

	obj.ReferenceCount++
	if rc, e := referenceMap[obj.ID]; e {
		referenceMap[obj.ID] = rc + 1
	} else {
		referenceMap[obj.ID] = 1
	}

	if tp.IsPrimitive {
		if requestMgr != nil {
			// 生成输入对象，存储到 inputMap ，并返回
			return obj, requestObj, nil
		}

		return obj, nil, nil
		//} else if tp.IsRawSlice {
		//	// 是数组
	} else if tp.IsStruct {
		// 是结构体
		nf := tp.Reference.NumField()
		lo := reflects.ParseStructLocation(tp.Reference)
		for n := 0; n < nf; n++ {
			fd := tp.Reference.Field(n)
			// 解析 tag，并获取描述信息
			fdTag := reflects.ParseTag(&fd)
			if fdTag == nil {
				// 忽略此字段
				continue
			}

			fdInfo := reflects.ParseField(lo, &fd)
			// 递归
			fdObj, _, err := traversal(lo, fdInfo.Raw, referenceMap, objMap, requestMgr)
			if err != nil {
				return nil, nil, err
			}

			// 所属对象不包含 Children 字段
			itemObj := &apimap.ObjectInfo{
				ReferenceID: fdObj.ID,
				Name:        fdTag.FieldName,
				Type:        fdObj.Type,
				Array:       fdInfo.IsRawSlice,
				Require:     fdTag.Require,
				Description: fdTag.Description,
				Deprecated:  fdTag.Deprecated,
			}
			obj.Children = append(obj.Children, itemObj)

			if requestMgr != nil {
				reqField := object.GenerateRequestObjectField(fdTag.FieldName, fd.Name, &fdInfo.TypeInfo, fdTag.Require)
				//reqField := &object.RequestObjectField{
				//	Name:         fdTag.FieldName,
				//	Type:         fdInfo.Reference,
				//	Ptr:          fdInfo.IsRawPtr,
				//	Slice:        fdInfo.IsRawSlice,
				//	SliceItemPtr: fdInfo.IsSliceItemIsPtr,
				//}
				requestObj.Fields = append(requestObj.Fields, reqField)
			}
		}

		//if requestMgr != nil {
		//	// 生成输入对象，存储到 inputMap ，并返回
		//}
		return obj, requestObj, nil
	}

	return nil, nil, nil
}

// decReferenceCount 减少引用
func decReferenceCount(referenceMap map[string]int, objMap map[string]*apimap.ObjectInfo) {
	for id, r := range referenceMap {
		if obj, e := objMap[id]; e {
			obj.ReferenceCount -= r
		}
	}
}
package objtraver

import (
	"fmt"
	"reflect"

	"github.com/seerx/runjson/pkg/context"

	"github.com/seerx/runjson/pkg/graph"

	"github.com/seerx/runjson/internal/reflects"

	"github.com/seerx/runjson/internal/runner/arguments/request"
)

// Traversal 遍历 typ
// 如果是遍历请求参数，则 inputMap 必须不为 nil ，并且会返回得到的 InputObject
// 如果是便器范库参数，则 inputMap 必须为 nil
// 如果 返回的 error 不为 nil，则说明参数有不合法的内容
func Traversal(loc *reflects.Location,
	typ reflect.Type,
	referenceMap map[string]int,
	objMap map[string]*graph.ObjectInfo,
	requestMgr *request.ObjectManager,
	log context.Log) (*graph.ObjectInfo, *request.RequestObject, error) {

	tp, err := reflects.ParseType(loc, typ)
	if err != nil {
		return nil, nil, err
	}
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
		refObj := &graph.ObjectInfo{
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
	obj = &graph.ObjectInfo{
		ID:       tp.ID(),
		Type:     tp.TypeName(),
		Children: nil,
	}
	objMap[obj.ID] = obj
	var requestObj *request.RequestObject
	if requestMgr != nil {
		// 生成 RequestObject 顶级对象
		requestObj = &request.RequestObject{
			//ID:        tp.ID(),
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

			fdInfo, err := reflects.ParseField(lo, &fd)
			if err != nil {
				return nil, nil, err
			}
			// 递归
			fdObj, _, err := Traversal(lo, fdInfo.Raw, referenceMap, objMap, requestMgr, log)
			if err != nil {
				return nil, nil, err
			}

			// 所属对象不包含 Children 字段
			itemObj := &graph.ObjectInfo{
				ReferenceID: fdObj.ID,
				Name:        fdTag.FieldName,
				Type:        fdObj.Type,
				Array:       fdInfo.IsRawSlice,
				Require:     fdTag.Require,
				Description: fdTag.Description,
				Range:       fdTag.Limit,
				Pattern:     fdTag.Regexp,
				Deprecated:  fdTag.Deprecated,
			}
			obj.Children = append(obj.Children, itemObj)

			if requestMgr != nil {
				reqField := request.GenerateRequestObjectField(fdTag, fd.Name, &fdInfo.TypeInfo, fdTag.Require, func(err error) {
					if err != nil {
						log.Warn("%s.%s: %s", lo.String(), fd.Name, err.Error())
					}
				})
				//reqField := &request.RequestObjectField{
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

// DecReferenceCount 减少引用
func DecReferenceCount(referenceMap map[string]int, objMap map[string]*graph.ObjectInfo) {
	for id, r := range referenceMap {
		if obj, e := objMap[id]; e {
			obj.ReferenceCount -= r
		}
	}
}

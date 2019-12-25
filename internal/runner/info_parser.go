package runner

import (
	"reflect"

	"github.com/seerx/runjson/internal/types"

	"github.com/seerx/runjson/pkg/intf"
)

func TryToParseFuncInfo(loader interface{}, loaderType reflect.Type, funcName string) *intf.FuncInfo {
	m, exists := loaderType.MethodByName(funcName + "Info")
	if !exists {
		return nil
	}

	mType := m.Type
	if mType.NumIn() != 1 || mType.NumOut() != 1 {
		return nil
	}

	out := m.Func.Call([]reflect.Value{reflect.ValueOf(loader)})

	outType := mType.Out(0)
	//if outType.Kind() == reflect.String {
	//	// 只有说明信息
	//	if desc, ok := out[0].Interface().(string); ok {
	//		return &intf.FuncInfo{
	//			Description: desc,
	//		}
	//	}
	//} else
	if yes, ptr := types.IsFuncInfo(outType); yes {
		if ptr {
			if info, ok := out[0].Interface().(*intf.FuncInfo); ok {
				return info
			}
		} else {
			if info, ok := out[0].Interface().(intf.FuncInfo); ok {
				return &info
			}
		}
	}
	return nil
}

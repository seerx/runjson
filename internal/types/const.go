package types

import (
	"reflect"

	"github.com/seerx/runjson/pkg/intf"
)

var (
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	requireType     = reflect.TypeOf((*intf.Require)(nil)).Elem()
	injectParamType = reflect.TypeOf((*map[string]interface{})(nil)).Elem()
)

func IsInjectParam(typ reflect.Type) bool {
	return injectParamType == typ
}

func IsRequirement(typ reflect.Type) bool {
	//if typ.Kind() == reflect.Ptr {
	//	return typ.Elem() == requireType, true
	//}
	//return typ == requireType, false
	return typ == requireType
}

// IsError 是不是错误接口
func IsError(typ reflect.Type) bool {
	return errorType == typ
}

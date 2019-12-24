package types

import "reflect"

var (
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	injectParamType = reflect.TypeOf((*map[string]interface{})(nil)).Elem()
)

func IsInjectParam(typ reflect.Type) bool {
	return injectParamType == typ
}

// IsError 是不是错误接口
func IsError(typ reflect.Type) bool {
	return errorType == typ
}

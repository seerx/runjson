package types

import (
	"reflect"

	"github.com/seerx/runjson/pkg/intf"
)

var (
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	fieldMapType    = reflect.TypeOf((*intf.FieldMap)(nil)).Elem()
	injectParamType = reflect.TypeOf((*map[string]interface{})(nil)).Elem()
)

func IsInjectParam(typ reflect.Type) bool {
	return injectParamType == typ
}

func IsFieldMap(typ reflect.Type) (bool, bool) {
	if typ.Kind() == reflect.Ptr {
		return typ.Elem() == fieldMapType, true
	}
	return typ == fieldMapType, false
}

// IsError 是不是错误接口
func IsError(typ reflect.Type) bool {
	return errorType == typ
}

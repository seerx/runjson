package types

import (
	"reflect"

	"github.com/seerx/runjson/pkg/intf"
)

var (
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	requireType     = reflect.TypeOf((*intf.Require)(nil)).Elem()
	funcInfoType    = reflect.TypeOf((*intf.FuncInfo)(nil)).Elem()
	injectParamType = reflect.TypeOf((*map[string]interface{})(nil)).Elem()
	resultsType     = reflect.TypeOf((*intf.Results)(nil)).Elem()
)

func IsInjectParam(typ reflect.Type) bool {
	return injectParamType == typ
}

func IsResults(typ reflect.Type) bool {
	//if typ.Kind() == reflect.Ptr {
	//	return typ.Elem() == responsesType, true
	//}
	return typ == resultsType
}

func IsFuncInfo(typ reflect.Type) (bool, bool) {
	if typ.Kind() == reflect.Ptr {
		return typ.Elem() == funcInfoType, true
	}
	return typ == funcInfoType, false
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

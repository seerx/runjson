package types

import (
	"reflect"

	"github.com/seerx/runjson/pkg/rj"
)

var (
	// NilType nil 类型
	NilType         = reflect.ValueOf(nil)
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	requireType     = reflect.TypeOf((*rj.Require)(nil)).Elem()
	funcInfoType    = reflect.TypeOf((*rj.FuncInfo)(nil)).Elem()
	injectParamType = reflect.TypeOf((*rj.InjectArg)(nil)).Elem()
	// injectParamType = reflect.TypeOf((*map[string]interface{})(nil)).Elem()
	resultsType = reflect.TypeOf((*rj.ResponseContext)(nil)).Elem()
)

// IsInjectParam 是否是注入函数的参数
func IsInjectParam(typ reflect.Type) bool {
	return injectParamType == typ
}

// IsResults 是否是返回类型
func IsResults(typ reflect.Type) bool {
	//if typ.Kind() == reflect.Ptr {
	//	return typ.Elem() == responsesType, true
	//}
	return typ == resultsType
}

// IsFuncInfo 是否是函数说明类型
func IsFuncInfo(typ reflect.Type) (bool, bool) {
	if typ.Kind() == reflect.Ptr {
		return typ.Elem() == funcInfoType, true
	}
	return typ == funcInfoType, false
}

// IsRequirement 是否是 Requires 类型
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

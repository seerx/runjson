package reflects

import (
	"reflect"
	"strings"
)

// ReflectStructType 获取结构类型的真实 Type
func ReflectStructType(obj interface{}) reflect.Type {
	typ := reflect.TypeOf(obj)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ
}

// IsDescFunc 是否描述函数
// 描述函数特点，以 Desc 为后缀，且返回一个 string
func IsDescFunc(method reflect.Method) bool {
	mn := method.Name
	if strings.HasSuffix(mn, "Desc") {
		typ := method.Type
		if typ.NumIn() == 1 && typ.NumOut() == 1 {
			out := typ.Out(0)
			return out.Kind() == reflect.String
		}
	}
	return false
}

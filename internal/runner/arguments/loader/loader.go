package loader

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/seerx/runjson/internal/types"

	"github.com/seerx/runjson/internal/runner/arguments"

	"github.com/seerx/runjson/internal/runner/inject"
)

var loaderMap = map[reflect.Type]*arguments.LoaderScheme{}

// ParseLoader 解析 Loader 结构
func ParseLoader(loaderType reflect.Type, injectorManager *inject.InjectorManager) *arguments.LoaderScheme {
	if ls, e := loaderMap[loaderType]; e {
		return ls
	}
	if loaderType.Kind() == reflect.Ptr {
		loaderType = loaderType.Elem()
	}
	ls := &arguments.LoaderScheme{Type: loaderType}

	for n := 0; n < loaderType.NumField(); n++ {
		field := loaderType.Field(n)

		if types.IsRequirement(field.Type) {
			if field.Name[:1] != strings.ToUpper(field.Name[:1]) {
				fmt.Printf("Field %s will not be set to require, because it's not export\n", field.Name)
				continue
			}
			ls.RequireFields = append(ls.RequireFields, field.Name)
			continue
		}
		ptr := false
		typ := field.Type
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
			ptr = true
		}
		if typ.Kind() == reflect.Struct || typ.Kind() == reflect.Interface {
			inj := injectorManager.Find(typ)
			if inj != nil {
				if field.Name[:1] != strings.ToUpper(field.Name[:1]) {
					fmt.Printf("Field %s will not inject, because it's not export\n", field.Name)
					continue
				}
				// 是注入类型
				ls.InjectFields = append(ls.InjectFields, &arguments.InjectField{
					Field:      field.Name,
					Injector:   inj,
					ValueIsPtr: ptr,
				})
			}
		}
		//reflects.ParseField()
	}

	loaderMap[loaderType] = ls
	return ls
}

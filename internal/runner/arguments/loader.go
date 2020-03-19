package arguments

import (
	"reflect"

	"github.com/seerx/runjson/internal/runner/inject"
)

type InjectField struct {
	Field      string
	Injector   *inject.Injector
	ValueIsPtr bool
}

// FieldResults 用于在 loader 中加载
type FieldResults struct {
	//ValueIsPtr bool
	FieldName string
}

// LoaderScheme loader 结构定义
type LoaderScheme struct {
	Type reflect.Type

	ResponsesFields []*FieldResults
	RequireFields   []string
	InjectFields    []*InjectField
}

func (ls *LoaderScheme) IsInject() bool {
	return false
}

func (ls *LoaderScheme) CreateValue(ctx *ArgumentContext) (reflect.Value, error) {
	inst := reflect.New(ls.Type)
	elem := inst.Elem()
	for _, fd := range ls.RequireFields {
		field := elem.FieldByName(fd)
		field.Set(reflect.ValueOf(ctx.Requirement))
	}

	for _, fd := range ls.ResponsesFields {
		field := elem.FieldByName(fd.FieldName)
		field.Set(reflect.ValueOf(ctx.Results))
	}

	for _, fd := range ls.InjectFields {
		field := elem.FieldByName(fd.Field)
		var val reflect.Value
		var exists bool
		if val, exists = ctx.InjectValueMap[fd.Injector.Type]; !exists {
			var err error
			val, err = fd.Injector.Call(ctx.ServiceName, ctx.Results, ctx.Param)
			if err != nil {
				panic(err)
			}
			ctx.InjectValueMap[fd.Injector.Type] = val
		}

		if fd.ValueIsPtr || fd.Injector.ReturnTypeIsInterface {
			field.Set(val)
		} else {
			field.Set(val.Elem())
		}
	}

	return inst, nil
}

//func (ls *LoaderScheme) IsInjectInterface() bool {
//	return false
//}

//type LoaderScheme struct {
//}
//
//func (ls *LoaderScheme) New() *LoaderScheme {
//	return &LoaderScheme{}
//}

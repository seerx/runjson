package arguments

import (
	"reflect"

	"github.com/seerx/runjson/internal/runner/inject"
)

type ArgInjector struct {
	Injector    *inject.Injector
	ValueIsPtr  bool
	IsInterface bool
}

func (a *ArgInjector) IsInject() bool {
	return true
}

func (a *ArgInjector) CreateValue(ctx *ArgumentContext) reflect.Value {
	var val reflect.Value
	var exists bool
	if val, exists = ctx.InjectValueMap[a.Injector.Type]; !exists {
		var err error
		val, err = a.Injector.Call(ctx.Param)
		if err != nil {
			panic(err)
		}
		ctx.InjectValueMap[a.Injector.Type] = val
	}

	if a.ValueIsPtr || a.Injector.ReturnTypeIsInterface {
		return val
	}
	return val.Elem()
}

//func (a *ArgInjector) IsInjectInterface() bool {
//	return a.IsInterface
//}

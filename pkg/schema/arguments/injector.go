package arguments

import (
	"reflect"

	"github.com/seerx/chain/pkg/inject"
)

type ArgInjector struct {
	Injector    *inject.Injector
	IsInterface bool
}

func (a *ArgInjector) CreateValue(ctx *ArgumentContext) reflect.Value {
	if val, e := ctx.InjectValueMap[a.Injector.Type]; e {
		return val
	}
	val, err := a.Injector.Call(ctx.Param)
	if err != nil {
		panic(err)
	}
	ctx.InjectValueMap[a.Injector.Type] = val
	return val
}

func (a *ArgInjector) IsInjectInterface() bool {
	return a.IsInterface
}

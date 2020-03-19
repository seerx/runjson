package arguments

import (
	"reflect"

	"github.com/seerx/runjson/internal/runner/inject"
)

// ArgInjector 注入参数
type ArgInjector struct {
	Injector    *inject.Injector
	ValueIsPtr  bool
	IsInterface bool
}

// IsInject 是否注入函数
func (a *ArgInjector) IsInject() bool {
	return true
}

// CreateValue 创建值
func (a *ArgInjector) CreateValue(ctx *ArgumentContext) (reflect.Value, error) {
	var val reflect.Value
	var exists bool
	if val, exists = ctx.InjectValueMap[a.Injector.Type]; !exists {
		var err error
		val, err = a.Injector.Call(ctx.ServiceName, ctx.Results, ctx.Param)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
		ctx.InjectValueMap[a.Injector.Type] = val
	}

	if a.ValueIsPtr || a.Injector.ReturnTypeIsInterface {
		return val, nil
	}
	return val.Elem(), nil
}

//func (a *ArgInjector) IsInjectInterface() bool {
//	return a.IsInterface
//}

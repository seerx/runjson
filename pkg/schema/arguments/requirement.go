package arguments

import (
	"reflect"
)

type ArgRequire struct {
	IsPtr bool
}

func (a *ArgRequire) CreateValue(ctx *ArgumentContext) reflect.Value {
	if a.IsPtr {
		return reflect.ValueOf(ctx.FieldMap)
	}
	return reflect.ValueOf(ctx.FieldMap).Elem()
	//panic("implement me")
}

func (a *ArgRequire) IsInjectInterface() bool {
	return false
}

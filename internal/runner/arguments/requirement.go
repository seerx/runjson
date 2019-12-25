package arguments

import (
	"reflect"
)

type ArgRequire struct {
	//IsPtr bool
}

func (a *ArgRequire) CreateValue(ctx *ArgumentContext) reflect.Value {
	return reflect.ValueOf(ctx.Requirement)
}

func (a *ArgRequire) IsInjectInterface() bool {
	return false
}

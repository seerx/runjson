package arguments

import (
	"reflect"
)

type ArgRequire struct {
	//IsPtr bool
}

func (a *ArgRequire) IsInject() bool {
	return false
}

func (a *ArgRequire) CreateValue(ctx *ArgumentContext) (reflect.Value, error) {
	return reflect.ValueOf(ctx.Requirement), nil
}

//func (a *ArgRequire) IsInjectInterface() bool {
//	return false
//}

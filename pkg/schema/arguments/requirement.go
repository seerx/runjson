package arguments

import (
	"reflect"
)

type ArgRequire struct {
	IsPtr bool
}

func (a *ArgRequire) CreateValue(ctx *ArgumentContext) reflect.Value {

	panic("implement me")
}

func (a *ArgRequire) IsInjectInterface() bool {
	return false
}

package arguments

import (
	"reflect"

	"github.com/seerx/runjson/internal/object"
)

type ArgRequest struct {
	ArgField *object.RequestObjectField
	Arg      *object.RequestObject
}

func (a *ArgRequest) CreateValue(ctx *ArgumentContext) reflect.Value {
	if ctx.RequestArgument == nil {
		return reflect.ValueOf(nil)
	}
	//return *ctx.RequestArgument
	if a.ArgField.Ptr {
		return *ctx.RequestArgument
	}
	return (*ctx.RequestArgument).Elem()
}

func (a *ArgRequest) IsInjectInterface() bool {
	return false
}

package arguments

import (
	"reflect"

	"github.com/seerx/runjson/internal/runner/arguments/request"
)

type ArgRequest struct {
	ArgField *request.RequestObjectField
	Arg      *request.RequestObject
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

//func (a *ArgRequest) IsInjectInterface() bool {
//	return false
//}

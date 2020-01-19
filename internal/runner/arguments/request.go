package arguments

import (
	"reflect"

	"github.com/seerx/runjson/internal/runner/arguments/request"
)

type ArgRequest struct {
	ArgField *request.RequestObjectField
	Arg      *request.RequestObject
}

func (a *ArgRequest) IsInject() bool {
	return false
}

func (a *ArgRequest) CreateValue(ctx *ArgumentContext) (reflect.Value, error) {
	if ctx.RequestArgument == nil {
		return reflect.ValueOf(nil), nil
	}
	//return *ctx.RequestArgument
	if a.ArgField.Ptr {
		return *ctx.RequestArgument, nil
	}
	return (*ctx.RequestArgument).Elem(), nil
}

//func (a *ArgRequest) IsInjectInterface() bool {
//	return false
//}

package arguments

import (
	"reflect"

	"github.com/seerx/chain/internal/object"
)

type ArgRequest struct {
	Arg *object.RequestObject
}

func (a *ArgRequest) CreateValue(ctx *ArgumentContext) reflect.Value {
	return *ctx.RequestArgument
}

func (a *ArgRequest) IsInjectInterface() bool {
	return false
}

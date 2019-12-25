package arguments

import (
	"reflect"

	"github.com/seerx/runjson/pkg/intf"
)

// ArgumentContext 参数上下文
type ArgumentContext struct {
	Param           map[string]interface{}
	InjectValueMap  map[reflect.Type]reflect.Value
	RequestArgument *reflect.Value
	Requirement     intf.Require
}

// 参数接口
type Argument interface {
	CreateValue(ctx *ArgumentContext) reflect.Value
	IsInjectInterface() bool
}

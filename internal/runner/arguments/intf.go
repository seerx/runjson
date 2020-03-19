package arguments

import (
	"reflect"

	"github.com/seerx/runjson/pkg/rj"
)

// ArgumentContext 参数上下文
type ArgumentContext struct {
	Param           map[string]interface{}
	InjectValueMap  map[reflect.Type]reflect.Value
	RequestArgument *reflect.Value
	Requirement     rj.Require
	Results         rj.ResponseContext
	ServiceName     string
}

// 参数接口
type Argument interface {
	CreateValue(ctx *ArgumentContext) (reflect.Value, error)
	IsInject() bool
	//AsClearTask() *rj.OnComplete
	//IsInjectInterface() bool
}

package arguments

import "reflect"

// ArgumentContext 参数上下文
type ArgumentContext struct {
	Param           map[string]interface{}
	InjectValueMap  map[reflect.Type]reflect.Value
	RequestArgument *reflect.Value
}

// 参数接口
type Argument interface {
	CreateValue(ctx *ArgumentContext) reflect.Value
	IsInjectInterface() bool
}

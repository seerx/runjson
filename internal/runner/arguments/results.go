package arguments

import "reflect"

// ArgResults 前面执行过的返回值
type ArgResults struct {
	//ValueIsPtr bool
}

func (a *ArgResults) IsInject() bool {
	return false
}

func (a *ArgResults) CreateValue(ctx *ArgumentContext) reflect.Value {
	//if a.ValueIsPtr {
	return reflect.ValueOf(ctx.Results)
	//}
	//return reflect.ValueOf(ctx.Results).Elem()
}

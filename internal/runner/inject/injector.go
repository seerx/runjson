package inject

import (
	"reflect"

	"github.com/seerx/runjson/internal/reflects"
)

// Injector 注入结构
type Injector struct {
	Type                  reflect.Type  // 注入类型，即注入函数返回的类型
	Func                  reflect.Value // 注入函数实例
	ReturnTypeIsInterface bool
	Location              *reflects.Location
}

// Call 调用注入函数
func (i *Injector) Call(arg map[string]interface{}) (reflect.Value, error) {
	args := []reflect.Value{reflect.ValueOf(arg)}
	res := i.Func.Call(args)
	if res[1].IsNil() {
		return res[0], nil
	}

	return res[0], res[1].Interface().(error)
}

package inject

import (
	"reflect"

	"github.com/seerx/runjson/internal/reflects"
	"github.com/seerx/runjson/pkg/rj"
)

// Injector 注入结构
type Injector struct {
	Type                  reflect.Type  // 注入类型，即注入函数返回的类型
	Func                  reflect.Value // 注入函数实例
	ArgIsPtr              bool          // 参数是否指针
	ReturnTypeIsInterface bool
	Location              *reflects.Location
}

// Call 调用注入函数
func (i *Injector) Call(serviceName string, response rj.ResponseContext, arg map[string]interface{}) (reflect.Value, error) {
	injectArg := &rj.InjectArg{
		Service:  serviceName,
		Response: response,
		Args:     arg,
	}
	var args []reflect.Value
	if i.ArgIsPtr {
		args = []reflect.Value{reflect.ValueOf(injectArg)}
	} else {
		args = []reflect.Value{reflect.ValueOf(*injectArg)}
	}

	res := i.Func.Call(args)
	if res[1].IsNil() {
		return res[0], nil
	}

	return res[0], res[1].Interface().(error)
}

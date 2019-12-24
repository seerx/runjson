package schema

import (
	"fmt"
	"reflect"

	"github.com/seerx/chain/internal/object"

	"github.com/seerx/chain/pkg/context"

	"github.com/seerx/chain/pkg/inject"

	"github.com/seerx/chain/pkg/schema/arguments"

	"github.com/seerx/chain/internal/reflects"
)

// Service 服务定义
type Service struct {
	Name     string // 服务名称
	method   reflect.Method
	loader   reflect.Type       // 函数所属结构体类型，非指针
	funcType reflect.Type       // 函数 Type
	location *reflects.Location // 函数位置

	injectMgr *inject.InjectorManager

	requestObjectMgr *object.RequestObjectManager

	returnType     reflect.Type               // 函数有效返回值 Type
	returnObjectID string                     // 返回类型 ID
	requestObject  *object.RequestObjectField // 函数接收值的 Type
	inputArgs      []arguments.Argument       // 函数输入参数表

	loaderStruct *arguments.LoaderStruct
}

func (s *Service) Run(ctx *context.Context, argument interface{}) (interface{}, error) {
	var arg *reflect.Value
	if s.requestObject != nil {
		a, err := s.requestObject.NewInstance("", argument, s.requestObjectMgr)
		if err != nil {
			return nil, err
		}
		arg = &a
	}

	// 组织函数参数
	argContext := &arguments.ArgumentContext{
		Param:           ctx.Param,
		RequestArgument: arg,
		InjectValueMap:  map[reflect.Type]reflect.Value{},
	}

	args := make([]reflect.Value, len(s.inputArgs), len(s.inputArgs))
	for n, a := range s.inputArgs {
		argVal := a.CreateValue(argContext)
		// 判断是否实现 io.Closer 接口
		args[n] = argVal
	}

	// call 函数
	res := s.method.Func.Call(args)
	if res == nil || len(res) != 2 {
		// 没有返回值，或这返回值的数量不是两个
		return nil, fmt.Errorf("Resolver <%s> error, need return values", s.Name)
	}

	out := res[0].Interface()
	errOut := res[1].Interface()
	var err error = nil
	if errOut != nil {
		ok := false
		err, ok = errOut.(error)
		if !ok {
			return nil, fmt.Errorf("Resolver <%s> error, second return must be error", s.Name)
		}
	}

	return out, err
}

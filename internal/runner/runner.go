package runner

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/seerx/runjson/internal/types"

	"github.com/seerx/runjson/pkg/rj"

	"github.com/seerx/runjson/internal/runner/arguments"
	"github.com/seerx/runjson/internal/runner/inject"

	"github.com/seerx/runjson/internal/runner/arguments/fieldmap"

	"github.com/seerx/runjson/internal/runner/arguments/request"

	"github.com/seerx/runjson/pkg/context"

	"github.com/seerx/runjson/internal/reflects"
)

// Runners 服务信息，用于执行服务
type Runners struct {
	// 以服务名称为 key
	RunnerMap map[string]*JSONRunner
	// 以服务函数的 位置 为 key
	runnerMap map[string]*JSONRunner
}

// New 创建 Runners
func New() *Runners {
	return &Runners{
		RunnerMap: map[string]*JSONRunner{},
		runnerMap: map[string]*JSONRunner{},
	}
}

var regOfStructPointer = regexp.MustCompile("\\(\\*([\\w\\d_]+)\\)")

func parseFuncPath(fn reflect.Value) string {
	name := runtime.FuncForPC(fn.Pointer()).Name()
	// 去掉函数名称后面的 -fm
	dotPos := strings.LastIndex(name, ".")
	linePos := strings.LastIndex(name, "-")
	if linePos > dotPos {
		name = name[:linePos]
	}
	// 去掉指针引用方式  (*struct)  => struct
	//reg := regexp.MustCompile("\\(\\*([\\w\\d_]+)\\)")
	name = regOfStructPointer.ReplaceAllStringFunc(name, func(s string) string {
		return s[2 : len(s)-1]
	})

	return name
}

// Add 添加
func (r *Runners) Add(runner *JSONRunner) {
	r.RunnerMap[runner.Name] = runner
	//name := runtime.FuncForPC(runner.method.Func.Pointer()).Name()
	//fmt.Println(name)
	//loc := reflects.ParseStructFuncLocation(runner.loader, runner.method)
	name := parseFuncPath(runner.method.Func)
	//fmt.Println("Add:", name)
	r.runnerMap[name] = runner
}

// Find 根据函数查找 JSONRunner
func (r *Runners) Find(method interface{}) (*JSONRunner, error) {
	//name := runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name()
	//dotPos := strings.LastIndex(name, ".")
	//linePos := strings.LastIndex(name, "-")
	//if linePos > dotPos {
	//	name = name[:linePos]
	//}
	name := parseFuncPath(reflect.ValueOf(method))
	//fmt.Println("Find:", name)
	if runner, ok := r.runnerMap[name]; ok {
		return runner, nil
	}
	return nil, fmt.Errorf("Runner [%s] is not found", name)
}

// Get 根据名称获取 JSONRunner
func (r *Runners) Get(runnerName string) *JSONRunner {
	return r.RunnerMap[runnerName]
}

// JSONRunner 服务定义
type JSONRunner struct {
	Name     string // 服务名称
	method   reflect.Method
	loader   reflect.Type       // 函数所属结构体类型，非指针
	funcType reflect.Type       // 函数 Type
	location *reflects.Location // 函数位置

	injectMgr *inject.InjectorManager

	requestObjectMgr *request.ObjectManager

	ReturnType           reflect.Type                // 函数有效返回值 Type
	ReturnObjectID       string                      // 返回类型 ID
	ReturnObjectIsArray  bool                        // 返回的数据是否数组
	requestObject        *request.RequestObjectField // 函数接收值的 Type
	RequestObjectID      string
	RequestObjectIsArray bool
	inputArgs            []arguments.Argument // 函数输入参数表
	AccessControllers    []*inject.Injector   // 权限控制注入列表

	loaderStruct *arguments.LoaderScheme
}

// SetRequestArgRequire 设置必须参数
func (s *JSONRunner) SetRequestArgRequire(require bool) {
	if s.requestObject != nil {
		s.requestObject.Require = require
	}
}

// cloneInjectMap 复制 injectMap
func cloneInjectMap(injectMap map[reflect.Type]reflect.Value) map[reflect.Type]reflect.Value {
	res := map[reflect.Type]reflect.Value{}
	if injectMap != nil {
		for k, v := range injectMap {
			res[k] = v
		}
	}
	return res
}

// Run 运行 json
func (s *JSONRunner) Run(ctx *context.Context, argument interface{}, injectMap map[reflect.Type]reflect.Value, results rj.ResponseContext) (interface{}, error) {
	var arg *reflect.Value
	fm := &fieldmap.FieldMap{}
	if s.requestObject != nil {
		a, err := s.requestObject.NewInstance("", argument, s.requestObjectMgr, fm)
		if err != nil {
			return nil, err
		}
		arg = &a
	}
	if injectMap == nil {
		injectMap = map[reflect.Type]reflect.Value{}
	}
	// 组织函数参数
	argContext := &arguments.ArgumentContext{
		ServiceName:     s.Name,
		Param:           ctx.Param,
		RequestArgument: arg,
		InjectValueMap:  cloneInjectMap(injectMap),
		Requirement:     fm,
		Results:         results,
	}

	clearTasks := []rj.OnComplete{}
	defer func() {
		for _, task := range clearTasks {
			task.Clear()
		}
	}()

	args := make([]reflect.Value, len(s.inputArgs), len(s.inputArgs))
	for n, a := range s.inputArgs {
		argVal, err := a.CreateValue(argContext)
		if err != nil {
			return nil, err
		}
		if a.IsInject() {
			// 注入类型
			// 判断是否实现 io.Closer 接口
			if task := asClearTask(argVal); task != nil {
				clearTasks = append(clearTasks, task)
			}
		}
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

func asClearTask(value reflect.Value) rj.OnComplete {
	if value == types.NilType {
		return nil
	}
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return nil
		}
	}
	val := value.Interface()
	if val == nil {
		return nil
	}

	task, ok := val.(rj.OnComplete)
	if ok {
		return task
	}
	return nil
}

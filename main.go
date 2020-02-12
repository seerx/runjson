package runjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/seerx/runjson/internal/util"

	"github.com/seerx/runjson/internal/runner"

	"github.com/seerx/runjson/pkg/graph"

	"github.com/seerx/runjson/pkg/context"

	"github.com/seerx/runjson/internal/runner/arguments/request"

	"github.com/seerx/runjson/internal/runner/inject"

	"github.com/seerx/runjson/pkg/rj"
)

// Error RunJson 错误信息
type Error struct {
	Err     error
	ctx     *context.Context
	request rj.Requests
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return ""
}

type (
	BeforeRun func(*context.Context, rj.Requests) error
	AfterRun  func(*context.Context, rj.Requests, rj.ResponseContext) error

	BeforeExecute func(ctx *context.Context, item *rj.Request) error
	AfterExecute  func(ctx *context.Context, item *rj.Request, result *rj.ResponseItem, results rj.ResponseContext) error

	OnError func(err *Error)
)

// Runner 结构体
type Runner struct {
	// 用于对外接口文档
	ApiInfo *graph.ApiInfo
	// 用于执行服务
	service *runner.Runners
	// 日志
	log context.Log
	// 注册的信息
	loaders []rj.Loader
	// 注入管理
	injector *inject.InjectorManager
	// 请求参数管理
	requestObjectManager *request.RequestObjectManager
	//groups  []*runner.Group
	//funcs   map[string]*runner.JSONRunner
	beforeRun     BeforeRun
	beforeExecute BeforeExecute
	afterRun      AfterRun
	afterExecute  AfterExecute
	onError       OnError
}

// SetLogger 设置日志输出
func (r *Runner) SetLogger(log context.Log) *Runner {
	r.log = log
	return r
}

// ErrorHandler 错误处理函数
func (r *Runner) ErrorHandler(handler OnError) *Runner {
	r.onError = handler
	return r
}

func (r *Runner) BeforeRun(fn BeforeRun) *Runner {
	r.beforeRun = fn
	return r
}
func (r *Runner) BeforeExecute(fn BeforeExecute) *Runner {
	r.beforeExecute = fn
	return r
}
func (r *Runner) AfterRun(fn AfterRun) *Runner {
	r.afterRun = fn
	return r
}
func (r *Runner) AfterExecute(fn AfterExecute) *Runner {
	r.afterExecute = fn
	return r
}

type results struct {
	response rj.Response
	run      *Runner
	count    int
	index    int
}

func (r *results) CallCount() int {
	return r.count
}

func (r *results) CallIndex() int {
	return r.index
}

func (r *results) Get(method interface{}) ([]*rj.ResponseItem, error) {
	if jr, err := r.run.service.Find(method); err != nil {
		return nil, err
	} else {
		if rsp, exists := r.response[jr.Name]; !exists {
			return nil, fmt.Errorf("Result of [%s] not found", jr.Name)
		} else {
			return rsp, nil
		}
	}
}

func New() *Runner {
	//log := logrus.Logger{
	//	Level:     logrus.WarnLevel,
	//	Formatter: &logrus.TextFormatter{},
	//}
	return &Runner{
		ApiInfo: &graph.ApiInfo{
			Groups:   nil,
			Request:  map[string]*graph.ObjectInfo{},
			Response: map[string]*graph.ObjectInfo{},
		},
		log:                  &util.Logger{},
		loaders:              nil,
		service:              runner.New(),
		injector:             inject.NewManager(),
		requestObjectManager: request.NewRequestObjectManager(),
	}
}

// Register 注册功能
func (r *Runner) Register(loaders ...rj.Loader) {
	r.loaders = append(r.loaders, loaders...)
}

// Inject 注册注入函数
func (r *Runner) Inject(fns ...interface{}) error {
	for _, fn := range fns {
		if err := r.injector.Register(fn); err != nil {
			return err
		}
	}
	return nil
}

// InjectProxy 注册代理注入
func (r *Runner) InjectProxy(fn interface{}, injectType reflect.Type, proxyFn interface{}) error {
	return r.injector.RegisterWithProxy(fn, injectType, proxyFn)
}

func (r *Runner) execute(ctx *context.Context, request *rj.Request, rslt *results, onResponse func(key string, rsp *rj.ResponseItem)) {
	defer func() {
		if err := recover(); err != nil {
			onResponse(request.Service, &rj.ResponseItem{
				Error: fmt.Sprintf("%v", err),
				Data:  nil,
			})
		}
	}()
	//resKey := request.Service
	var rsp *rj.ResponseItem
	svc := r.service.Get(request.Service)
	if svc != nil {
		res, err := svc.Run(ctx, request.Args, rslt)

		if err != nil {
			rsp = &rj.ResponseItem{
				Error: err.Error(),
			}
		} else {
			rsp = &rj.ResponseItem{
				Error: "",
				Data:  res,
			}
		}
	} else {
		rsp = &rj.ResponseItem{
			Error: "No service named " + request.Service,
		}
	}

	onResponse(request.Service, rsp)
}

func (r *Runner) doRun(ctx *context.Context, reqs rj.Requests, returnFn func(rj.Response, error)) {
	defer func() {
		if err := recover(); err != nil {
			returnFn(nil, errors.New(err.(string)))
		}
	}()
	//r.log.Debug("Requests: \n%s", data)

	if r.beforeRun != nil {
		if err := r.beforeRun(ctx, reqs); err != nil {
			returnFn(nil, err)
			return
		}
	}

	response := rj.Response{}
	rslt := &results{
		response: response,
		run:      r,
		count:    len(reqs),
		index:    0,
	}

	for n, request := range reqs {
		// before
		if r.beforeExecute != nil {
			if err := r.beforeExecute(ctx, request); err != nil {
				returnFn(response, err)
				return
			}
		}
		var result *rj.ResponseItem
		rslt.index = n
		r.execute(ctx, request, rslt, func(key string, rsp *rj.ResponseItem) {
			if resAry, exists := response[request.Service]; exists {
				response[key] = append(resAry, rsp)
			} else {
				response[key] = []*rj.ResponseItem{rsp}
			}
			result = rsp
		})
		// after
		if r.afterExecute != nil {
			if err := r.afterExecute(ctx, request, result, rslt); err != nil {
				returnFn(response, err)
				return
			}
		}

		//r.log.Debug("Call: %s", request.Service)
	}

	if r.afterRun != nil {
		if err := r.afterRun(ctx, reqs, rslt); err != nil {
			returnFn(response, err)
			return
		}
	}

	returnFn(response, nil)
}

// RunString 运行字符串形式的参数
func (r *Runner) RunString(ctx *context.Context, data string) (rj.Response, error) {
	var rsp rj.Response
	var err error
	var reqs = rj.Requests{}
	err = json.Unmarshal([]byte(data), &reqs)
	if err != nil {
		r.log.Error(err, "json.Unmarshal")
		if r.onError != nil {
			r.onError(&Error{
				Err:     err,
				ctx:     ctx,
				request: reqs,
			})
		}
		return nil, err
	}
	r.doRun(ctx, reqs, func(responses rj.Response, e error) {
		rsp = responses
		err = e
		if r.onError != nil {
			r.onError(&Error{
				Err:     err,
				ctx:     ctx,
				request: reqs,
			})
		}
	})
	return rsp, err
}

// RunRequests 运行 rj.Requests 形式的参数
func (r *Runner) RunRequests(ctx *context.Context, reqs rj.Requests) (rj.Response, error) {
	var rsp rj.Response
	var err error
	r.doRun(ctx, reqs, func(responses rj.Response, e error) {
		rsp = responses
		err = e
		if r.onError != nil {
			r.onError(&Error{
				Err:     err,
				ctx:     ctx,
				request: reqs,
			})
		}
	})
	return rsp, err
}

// Engage 解析功能，以启动功能 from Star Trek
func (r *Runner) Engage() error {
	// 用于接口的 API 文档信息
	amap := r.ApiInfo

	for _, loader := range r.loaders {
		grpInfo := loader.Group()
		if grpInfo == nil {
			return errors.New("A loader must belong a group")
		}
		if strings.TrimSpace(grpInfo.Name) == "" {
			return errors.New("Group's name shuld not be empty")
		}
		// 获得分组
		grp := amap.GetGroup(grpInfo.Name, grpInfo.Description)
		// 解析 functions
		loaderTyp := reflect.TypeOf(loader)

		// 遍历函数
		nm := loaderTyp.NumMethod()
		for n := 0; n < nm; n++ {
			method := loaderTyp.Method(n)

			if rj.GROUP_FUNC == method.Name {
				// 跳过 Loader 接口的函数
				continue
			}

			// 生成 JSONRunner 名称
			svcName, err := grp.GenerateServiceName(method.Name)
			if err != nil {
				// 同名 service 已经存在
				r.log.Error(err, "JSONRunner exists")
				continue
			}
			// 解析服务函数
			svc, err := runner.TryParserAsService(loaderTyp,
				r.injector,
				r.requestObjectManager,
				method,
				amap,
				r.log)
			if err != nil {
				// 不是合法的服务函数
				//r.log.Warn("[%s] is not a service function: %v\n", svcName, err)
				continue
			}

			if svc != nil {
				// 解析完成
				// 填写服务名称
				svc.Name = svcName

				svcInfo := &graph.ServiceInfo{
					Name:           svcName,
					InputObjectID:  svc.RequestObjectID,
					InputIsArray:   svc.RequestObjectIsArray,
					OutputObjectID: svc.ReturnObjectID,
					OutputIsArray:  svc.ReturnObjectIsArray,
				}
				// 解析服务描述信息
				info := runner.TryToParseFuncInfo(loader, loaderTyp, method.Name)
				if info != nil {
					svcInfo.Description = info.Description
					svcInfo.Deprecated = info.Deprecated
					svcInfo.History = info.History
					svcInfo.InputIsRequire = info.InputIsRequire
					svc.SetRequestArgRequire(svcInfo.InputIsRequire)
				}
				grp.AddService(svcInfo)
				//grp. = append(grp.Funcs, fn)
				r.service.Add(svc)
			}
		}
	}

	return nil
}

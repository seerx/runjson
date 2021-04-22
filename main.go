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
	// BeforeRun 运行前拦截函数定义
	BeforeRun func(*context.Context, rj.Requests) error
	// AfterRun 运行后拦截函数定义
	AfterRun func(*context.Context, rj.Requests, rj.ResponseContext) error
	// BeforeExecute 执行单个 API 前拦截函数定
	BeforeExecute func(ctx *context.Context, item *rj.Request) error
	// AfterExecute 执行单个 API 后拦截函数定
	AfterExecute func(ctx *context.Context, item *rj.Request, result *rj.ResponseItem, results rj.ResponseContext) error
	// OnError 错误报告函数定义
	OnError func(err *Error)
)

// Runner 结构体
type Runner struct {
	// 用于对外接口文档
	APIInfo *graph.APIInfo
	// 用于执行服务
	service *runner.Runners
	// 日志
	log context.Log
	// logRegister 记录 api 注册过程
	logRegister bool
	// 注册的信息
	loaders []rj.Loader
	// 注入管理
	injector *inject.InjectorManager
	// 请求参数管理
	requestObjectManager *request.ObjectManager
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

// SetLogRegister 设置时候记录注册 api 过程
func (r *Runner) SetLogRegister(log bool) *Runner {
	r.logRegister = log
	return r
}

// ErrorHandler 错误处理函数
func (r *Runner) ErrorHandler(handler OnError) *Runner {
	r.onError = handler
	return r
}

// BeforeRun 在批量执行前拦截
func (r *Runner) BeforeRun(fn BeforeRun) *Runner {
	r.beforeRun = fn
	return r
}

// BeforeExecute 在单个任务执行时拦截
func (r *Runner) BeforeExecute(fn BeforeExecute) *Runner {
	r.beforeExecute = fn
	return r
}

// AfterRun 在批量执行后执行
func (r *Runner) AfterRun(fn AfterRun) *Runner {
	r.afterRun = fn
	return r
}

// AfterExecute 在单个任务执行后拦截
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

// CallCount 调用个数
func (r *results) CallCount() int {
	return r.count
}

// CallIndex 调用次序
func (r *results) CallIndex() int {
	return r.index
}

func (r *results) Get(method interface{}) ([]*rj.ResponseItem, error) {
	jr, err := r.run.service.Find(method)
	if err != nil {
		return nil, err
	}
	rsp, exists := r.response[jr.Name]
	if !exists {
		return nil, fmt.Errorf("result of [%s] not found", jr.Name)
	}
	return rsp, nil
}

// New 新建 Runner
func New() *Runner {
	//log := logrus.Logger{
	//	Level:     logrus.WarnLevel,
	//	Formatter: &logrus.TextFormatter{},
	//}
	return &Runner{
		APIInfo: &graph.APIInfo{
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

// RegisterProvider 注册注入函数
func (r *Runner) RegisterProvider(fns ...interface{}) error {
	for _, fn := range fns {
		if err := r.injector.Register(fn); err != nil {
			return err
		}
	}
	return nil
}

// RegisterAccessController 注册兼顾权限控制的注入函数
func (r *Runner) RegisterAccessController(fn interface{}) error {
	return r.injector.RegisterAccessController(fn)
}

// InjectProxy 注册代理注入
// func (r *Runner) InjectProxy(fn interface{}, injectType reflect.Type, proxyFn interface{}) error {
// 	return r.injector.RegisterWithProxy(fn, injectType, proxyFn)
// }

func (r *Runner) execute(ctx *context.Context, injectMap map[reflect.Type]reflect.Value, request *rj.Request, rslt *results, onResponse func(key string, rsp *rj.ResponseItem)) {
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
		res, err := svc.Run(ctx, request.Args, injectMap, rslt)
		if err != nil {
			rsp = &rj.ResponseItem{
				Error:    err.Error(),
				DataType: svc.ReturnType,
			}
		} else {
			rsp = &rj.ResponseItem{
				Error:    "",
				Data:     res,
				DataType: svc.ReturnType,
			}
		}
	} else {
		rsp = &rj.ResponseItem{
			Error: "No service named " + request.Service,
		}
	}

	onResponse(request.Service, rsp)
}

func (r *Runner) checkAccess(reqs rj.Requests, ctx *context.Context, responseContext rj.ResponseContext) (map[reflect.Type]reflect.Value, error) {
	accessInject := map[reflect.Type]reflect.Value{}
	for _, req := range reqs {
		svc := r.service.Get(req.Service)
		if svc == nil {
			// 找不到服务
			return nil, fmt.Errorf("no service named %s", req.Service)
		}
		for _, ac := range svc.AccessControllers {
			val, err := ac.Call(req.Service, responseContext, ctx.Param)
			if err != nil {
				return nil, err
			}
			accessInject[ac.Type] = val
		}
	}
	return accessInject, nil
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

	// 检查权限
	injectMap, err := r.checkAccess(reqs, ctx, rslt)
	if err != nil {
		returnFn(nil, err)
		return
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
		r.execute(ctx, injectMap, request, rslt, func(key string, rsp *rj.ResponseItem) {
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
	amap := r.APIInfo

	for _, loader := range r.loaders {
		grpInfo := loader.Group()
		if grpInfo == nil {
			return errors.New("a loader must belong a group")
		}
		if strings.TrimSpace(grpInfo.Name) == "" {
			return errors.New("group's name shuld not be empty")
		}
		// 获得分组
		grp := amap.GetGroup(grpInfo.Name, grpInfo.Description)
		// 解析 functions
		loaderTyp := reflect.TypeOf(loader)

		// 遍历函数
		nm := loaderTyp.NumMethod()
		for n := 0; n < nm; n++ {
			method := loaderTyp.Method(n)

			if rj.GroupFunc == method.Name {
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
			if r.logRegister {
				r.log.Info("Try to register api %s ...", svcName)
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
				if r.logRegister {
					r.log.Warn("[%s] is not a service function: %v\n", svcName, err)
				}

				// r.log.Warn("[%s] is not a service function: %v\n", svcName, err)
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

package runjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/seerx/runjson/internal/runner"

	"github.com/seerx/runjson/pkg/graph"

	"github.com/seerx/runjson/pkg/context"

	"github.com/seerx/runjson/internal/runner/arguments/request"

	"github.com/seerx/runjson/internal/runner/inject"

	"github.com/seerx/runjson/pkg/intf"

	"github.com/sirupsen/logrus"
)

type (
	BeforeRun func(*context.Context, intf.Requests)
	AfterRun  func(*context.Context, intf.Requests, intf.Results)

	BeforeExecute func(ctx *context.Context, item *intf.Request)
	AfterExecute  func(ctx *context.Context, item *intf.Request, result *intf.ResponseItem, results intf.Results)
)

// Runner 结构体
type Runner struct {
	// 用于对外接口文档
	ApiInfo *graph.ApiInfo
	// 用于执行服务
	service *runner.Runners
	// 日志
	log logrus.Logger
	// 注册的信息
	loaders []intf.Loader
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
	response intf.Response
	run      *Runner
}

func (r *results) Get(method interface{}) ([]*intf.ResponseItem, error) {
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
	log := logrus.Logger{
		Level:     logrus.WarnLevel,
		Formatter: &logrus.TextFormatter{},
	}
	return &Runner{
		ApiInfo: &graph.ApiInfo{
			Groups:   nil,
			Request:  map[string]*graph.ObjectInfo{},
			Response: map[string]*graph.ObjectInfo{},
		},
		log:                  log,
		loaders:              nil,
		service:              runner.New(),
		injector:             inject.NewManager(),
		requestObjectManager: request.NewRequestObjectManager(),
	}
}

// Register 注册功能
func (r *Runner) Register(loaders ...intf.Loader) {
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

func (r *Runner) execute(ctx *context.Context, request *intf.Request, rslt *results, onResponse func(key string, rsp *intf.ResponseItem)) {
	defer func() {
		if err := recover(); err != nil {
			onResponse(request.Service, &intf.ResponseItem{
				Error: err.(string),
				Data:  nil,
			})
		}
	}()
	//resKey := request.Service
	var rsp *intf.ResponseItem
	svc := r.service.Get(request.Service)
	if svc != nil {
		res, err := svc.Run(ctx, request.Args, rslt)

		if err != nil {
			rsp = &intf.ResponseItem{
				Error: err.Error(),
			}
		} else {
			rsp = &intf.ResponseItem{
				Error: "",
				Data:  res,
			}
		}
	} else {
		rsp = &intf.ResponseItem{
			Error: "No runner named " + request.Service,
		}
	}

	onResponse(request.Service, rsp)
}

func (r *Runner) doRun(ctx *context.Context, data string, returnFn func(intf.Response, error)) {
	defer func() {
		if err := recover(); err != nil {
			returnFn(nil, errors.New(err.(string)))
		}
	}()
	//r.log.Debug("Requests: \n%s", data)
	var reqs = intf.Requests{}
	err := json.Unmarshal([]byte(data), &reqs)
	if err != nil {
		r.log.WithError(err).Error("json.Unmarshal")
		returnFn(nil, err)
	}

	if r.beforeRun != nil {
		r.beforeRun(ctx, reqs)
	}

	response := intf.Response{}
	rslt := &results{
		response: response,
		run:      r,
	}

	for _, request := range reqs {
		// before
		if r.beforeExecute != nil {
			r.beforeExecute(ctx, request)
		}
		var result *intf.ResponseItem
		r.execute(ctx, request, rslt, func(key string, rsp *intf.ResponseItem) {
			if resAry, exists := response[request.Service]; exists {
				response[key] = append(resAry, rsp)
			} else {
				response[key] = []*intf.ResponseItem{rsp}
			}
			result = rsp
		})
		// after
		if r.afterExecute != nil {
			r.afterExecute(ctx, request, result, rslt)
		}

		//r.log.Debug("Call: %s", request.Service)
	}

	if r.afterRun != nil {
		r.afterRun(ctx, reqs, rslt)
	}

	returnFn(response, nil)
}

// Run 执行
func (r *Runner) Run(ctx *context.Context, data string) (intf.Response, error) {
	var rsp intf.Response
	var err error
	r.doRun(ctx, data, func(responses intf.Response, e error) {
		rsp = responses
		err = e
	})
	return rsp, err
	//var reqs = intf.Requests{}
	//err := json.Unmarshal([]byte(data), &reqs)
	//if err != nil {
	//	r.log.WithError(err).Error("json.Unmarshal")
	//	return nil, err
	//}
	//
	//if r.beforeRun != nil {
	//	r.beforeRun(reqs)
	//}
	//
	//response := intf.Response{}
	//rslt := &results{
	//	response: response,
	//	run:      r,
	//}
	//
	//for _, request := range reqs {
	//	// before
	//	if r.beforeExecute != nil {
	//		r.beforeExecute(request)
	//	}
	//	var result *intf.ResponseItem
	//	r.execute(ctx, request, rslt, func(key string, rsp *intf.ResponseItem) {
	//		if resAry, exists := response[request.Service]; exists {
	//			response[key] = append(resAry, rsp)
	//		} else {
	//			response[key] = []*intf.ResponseItem{rsp}
	//		}
	//		result = rsp
	//	})
	//	// after
	//	if r.afterExecute != nil {
	//		r.afterExecute(request, result, rslt)
	//	}
	//
	//	//r.log.Debug("Call: %s", request.Service)
	//}
	//
	//if r.afterRun != nil {
	//	r.afterRun(reqs, rslt)
	//}
	//
	//return response, nil
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

			if intf.GROUP_FUNC == method.Name {
				// 跳过 Loader 接口的函数
				continue
			}

			// 生成 JSONRunner 名称
			svcName, err := grp.GenerateServiceName(method.Name)
			if err != nil {
				// 同名 service 已经存在
				r.log.WithError(err).Error("JSONRunner exists")
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
				r.log.WithError(err).
					Debugf("[%s] is not a service function", svcName)
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
				}
				grp.AddService(svcInfo)
				//grp. = append(grp.Funcs, fn)
				r.service.Add(svc)
			}
		}
	}

	return nil
}

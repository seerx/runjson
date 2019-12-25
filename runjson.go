package runjson

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/seerx/runjson/internal/runner"

	"github.com/seerx/runjson/pkg/graph"

	"github.com/seerx/runjson/pkg/context"

	"github.com/seerx/runjson/internal/object"

	"github.com/seerx/runjson/internal/runner/inject"

	"github.com/seerx/runjson/pkg/intf"

	"github.com/sirupsen/logrus"
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
	requestObjectManager *object.RequestObjectManager
	//groups  []*runner.Group
	//funcs   map[string]*runner.JSONRunner
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
		log:     log,
		loaders: nil,
		service: &runner.Runners{
			RunnerMap: map[string]*runner.JSONRunner{},
		},
		injector:             inject.NewManager(),
		requestObjectManager: object.NewRequestObjectManager(),
	}
}

// Register 注册功能
func (c *Runner) Register(loaders ...intf.Loader) {
	c.loaders = append(c.loaders, loaders...)
}

// Inject 注册注入函数
func (c *Runner) Inject(fns ...interface{}) error {
	for _, fn := range fns {
		if err := c.injector.Register(fn); err != nil {
			return err
		}
	}
	return nil
}

// Execute 执行
func (c *Runner) Execute(ctx *context.Context, data string) (Responses, error) {
	c.log.Debug("Requests: \n%s", data)
	var reqs = Requests{}
	err := json.Unmarshal([]byte(data), &reqs)
	if err != nil {
		c.log.WithError(err).Error("json.Unmarshal")
		return nil, err
	}

	rsp := map[string]*Response{}

	for _, request := range reqs {
		resKey := request.Alias
		if resKey == "" {
			resKey = request.Service
		}
		c.log.Debug("Call: %s", request.Service)
		svc := c.service.Get(request.Service)
		if svc != nil {
			res, err := svc.Run(ctx, request.Args)

			if err != nil {
				rsp[resKey] = &Response{
					Error: err.Error(),
				}
			} else {
				rsp[resKey] = &Response{
					Error: "",
					Data:  res,
				}
			}
		} else {
			rsp[resKey] = &Response{
				Error: "No JSONRunner Named " + request.Service,
			}
		}
		//request.Method
	}
	return rsp, nil
}

// Explain 解释定义
func (c *Runner) Explain() error {
	// 用于接口的 API 文档信息
	amap := c.ApiInfo

	for _, loader := range c.loaders {
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
				log.WithError(err).Error("JSONRunner exists")
				continue
			}
			// 解析服务函数
			svc, err := runner.TryParserAsService(loaderTyp,
				c.injector,
				c.requestObjectManager,
				method,
				amap,
				c.log)
			if err != nil {
				// 不是合法的服务函数
				log.WithError(err).
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
					OutputObjectID: svc.ReturnObjectID,
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
				c.service.Add(svc)
			}
		}
	}

	return nil
}

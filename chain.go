package chain

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/seerx/chain/pkg/context"

	"github.com/seerx/chain/internal/object"

	"github.com/seerx/chain/pkg/inject"

	"github.com/seerx/chain/pkg/intf"

	"github.com/seerx/chain/pkg/apimap"

	"github.com/seerx/chain/pkg/schema"

	"github.com/sirupsen/logrus"
)

// Chain 结构体
type Chain struct {
	// 用于对外接口文档
	ApiMap *apimap.MapInfo
	// 用于执行服务
	service *schema.Services
	// 日志
	log logrus.Logger
	// 注册的信息
	loaders []intf.Loader
	// 注入管理
	injector *inject.InjectorManager
	// 请求参数管理
	requestObjectManager *object.RequestObjectManager
	//groups  []*schema.Group
	//funcs   map[string]*schema.Service
}

func New() *Chain {
	return &Chain{
		ApiMap: &apimap.MapInfo{
			Groups:   nil,
			Request:  map[string]*apimap.ObjectInfo{},
			Response: map[string]*apimap.ObjectInfo{},
		},
		log:     logrus.Logger{},
		loaders: nil,
		service: &schema.Services{
			ServiceMap: map[string]*schema.Service{},
		},
		injector:             inject.NewManager(),
		requestObjectManager: object.NewRequestObjectManager(),
	}
}

// Register 注册功能
func (c *Chain) Register(loaders ...intf.Loader) {
	c.loaders = append(c.loaders, loaders...)
}

// Execute 执行
func (c *Chain) Execute(ctx *context.Context, data string) (Responses, error) {
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
		svc := c.service.GetService(request.Service)
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
				Error: "No Service Named " + request.Service,
			}
		}
		//request.Method
	}
	return rsp, nil
}

// Explain 解释定义
func (c *Chain) Explain() error {
	// 用于接口的 API 文档信息
	amap := c.ApiMap

	for _, loader := range c.loaders {
		grpInfo := loader.Group()
		if grpInfo == nil {
			return errors.New("A loader must belong a group")
		}
		if strings.TrimSpace(grpInfo.Name) == "" {
			return errors.New("Group's name shuld not be empty")
		}
		// 获得分组
		grp := amap.GetGroup(grpInfo)
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

			// 生成 Service 名称
			svcName, err := grp.GenerateServiceName(method.Name)
			if err != nil {
				// 同名 service 已经存在
				log.WithError(err).Error("Service exists")
				continue
			}
			// 解析服务函数
			svc, err := schema.TryParserAsService(loaderTyp,
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
				// 解析服务描述信息
				//grp. = append(grp.Funcs, fn)
				c.service.AddService(svc)
			}
		}
	}

	return nil
}

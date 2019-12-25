package schema

import (
	"fmt"
	"reflect"

	"github.com/seerx/runjson/internal/object"
	"github.com/seerx/runjson/internal/reflects"
	"github.com/seerx/runjson/internal/types"
	"github.com/seerx/runjson/pkg/apimap"
	"github.com/seerx/runjson/pkg/inject"
	"github.com/seerx/runjson/pkg/schema/arguments"
	"github.com/sirupsen/logrus"
)

// TryParserAsService 尝试解析函数为服务
func TryParserAsService(loader reflect.Type,
	injectManager *inject.InjectorManager,
	requestObjectManager *object.RequestObjectManager,
	method reflect.Method,
	info *apimap.MapInfo,
	log logrus.Logger) (*Service, error) {

	// 生成服务基础信息
	svc := &Service{
		//ID:       util.MD5(loader.PkgPath(), loader.Name(), method.Name),
		//Name:     svcName,
		location:         reflects.ParseStructFuncLocation(loader, method),
		method:           method,
		loader:           loader,
		funcType:         method.Type,
		injectMgr:        injectManager,
		requestObjectMgr: requestObjectManager,
	}

	rMap := map[string]int{}
	// 解析输出参数
	var outObj *apimap.ObjectInfo
	var err error
	if outObj, err = checkOutArguments(svc, rMap, info.Response, log); err != nil {
		decReferenceCount(rMap, info.Response)
		return nil, err
	}
	svc.returnObjectID = outObj.ID

	// 解析输入参数
	//var inInfo *apimap.ObjectInfo
	var inReq *object.RequestObjectField
	inMap := map[string]int{}
	_, inReq, err = checkInArguments(svc, requestObjectManager, inMap, info.Request, log)
	if err != nil {
		decReferenceCount(rMap, info.Response)
		decReferenceCount(inMap, info.Request)
		return nil, err
	}
	if inReq != nil {
		svc.requestObject = inReq
	}
	return svc, nil
}

// 检查函数的输入参数
func checkInArguments(svc *Service,
	requestObjectManager *object.RequestObjectManager,
	referenceMap map[string]int,
	objMap map[string]*apimap.ObjectInfo,
	log logrus.Logger) (*apimap.ObjectInfo, *object.RequestObjectField, error) {
	ic := svc.funcType.NumIn()
	var inInfo *apimap.ObjectInfo
	var inObj *object.RequestObject
	var inObjField *object.RequestObjectField
	svc.inputArgs = make([]arguments.Argument, ic, ic)
	for n := 0; n < ic; n++ {
		in := svc.funcType.In(n)
		if is, ptr := types.IsFieldMap(in); is {
			// 用于判断必填字段检查的参数
			svc.inputArgs[n] = &arguments.ArgRequire{IsPtr: ptr}
			continue
		}

		typeInfo := reflects.ParseType(svc.location, in)
		if n == 0 {
			// 结构体
			svc.loaderStruct = arguments.ParseLoader(in, svc.injectMgr)
			svc.inputArgs[n] = svc.loaderStruct
			continue
		}

		// 注入字段
		if typeInfo.IsInterface {
			// 接口类型，必须是注入字段
			injector := svc.injectMgr.Find(typeInfo.Raw)
			if injector == nil {
				return nil, nil, fmt.Errorf("No injector exists with type %s: %s", typeInfo.Name, svc.location)
			}
			svc.inputArgs[n] = &arguments.ArgInjector{
				Injector:    injector,
				IsInterface: true,
			}
			continue
		}
		if typeInfo.IsStruct {
			// 结构类型，必须定义为指针
			if !typeInfo.IsRawPtr {
				return nil, nil, fmt.Errorf("A struct argument must be a pointer %s: %s", typeInfo.Name, svc.location)
			}
			// 优先认定注入类型
			injector := svc.injectMgr.Find(typeInfo.Reference)
			if injector != nil {
				// 在注入类型中找到
				svc.inputArgs[n] = &arguments.ArgInjector{
					Injector:    injector,
					IsInterface: false,
				}
			} else {
				// 不是注入类型，那就一定是输入参数
				if svc.requestObject != nil {
					// 已经有输入参数了
					return nil, nil, fmt.Errorf("A service only has one request argument %s: %s", typeInfo.Name, svc.location)
				}

				var err error
				if inInfo, inObj, err = traversal(svc.location, typeInfo.Raw, referenceMap, objMap, requestObjectManager); err != nil {
					return inInfo, nil, fmt.Errorf("Service function invalid type: %s -> %s", err, svc.location)
				}

				inObjField = object.GenerateRequestObjectField(nil, "", typeInfo, false)

				svc.inputArgs[n] = &arguments.ArgRequest{
					Arg:      inObj,
					ArgField: inObjField,
				}
			}
			continue
		}
		if typeInfo.IsPrimitive {
			// 原生类型，必须是输入参数
			if svc.requestObject != nil {
				// 已经有输入参数了
				return nil, nil, fmt.Errorf("A service only has one request argument %s: %s", typeInfo.Name, svc.location)
			}
			var err error
			if inInfo, inObj, err = traversal(svc.location, typeInfo.Raw, referenceMap, objMap, requestObjectManager); err != nil {
				return inInfo, nil, fmt.Errorf("Service function invalid type: %s -> %s", err, svc.location)
			}
			inObjField = object.GenerateRequestObjectField(nil, "", typeInfo, false)
			svc.inputArgs[n] = &arguments.ArgRequest{
				Arg:      inObj,
				ArgField: inObjField,
			}
			//svc.requestObject = inObjField
			continue
		}
		return inInfo, nil, fmt.Errorf("Invalid service function : %s", svc.location)
	}

	return inInfo, inObjField, nil
}

// 检查函数的返回参数
func checkOutArguments(svc *Service, referenceMap map[string]int, objMap map[string]*apimap.ObjectInfo, log logrus.Logger) (*apimap.ObjectInfo, error) {
	funcLoc := svc.location.String()
	oc := svc.funcType.NumOut()
	if oc != 2 {
		log.Debug("Service function Must has return 2 values:", funcLoc)
		return nil, fmt.Errorf("Service function Must has return 2 values: %s", funcLoc)
	}

	o := svc.funcType.Out(1)
	if !types.IsError(o) {
		return nil, fmt.Errorf("Service function's second return argument must be error: %s", funcLoc)
	}
	//rMap := map[string]int{}
	o = svc.funcType.Out(0)
	if outObj, _, err := traversal(svc.location, o, referenceMap, objMap, nil); err != nil {
		return nil, fmt.Errorf("Service function invalid type: %s -> %s", err, funcLoc)
	} else {
		//decReferenceCount(referenceMap, info)
		return outObj, nil
	}
}

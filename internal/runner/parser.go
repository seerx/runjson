package runner

import (
	"fmt"
	"reflect"

	"github.com/seerx/runjson/pkg/context"

	"github.com/seerx/runjson/internal/runner/arguments/loader"

	"github.com/seerx/runjson/internal/runner/arguments"
	"github.com/seerx/runjson/internal/runner/inject"
	objtraver "github.com/seerx/runjson/internal/runner/objtraver"

	"github.com/seerx/runjson/pkg/graph"

	"github.com/seerx/runjson/internal/reflects"
	"github.com/seerx/runjson/internal/runner/arguments/request"
	"github.com/seerx/runjson/internal/types"
)

// TryParserAsService 尝试解析函数为服务
func TryParserAsService(loader reflect.Type,
	injectManager *inject.InjectorManager,
	requestObjectManager *request.ObjectManager,
	method reflect.Method,
	info *graph.APIInfo,
	log context.Log) (*JSONRunner, error) {

	// 生成服务基础信息
	svc := &JSONRunner{
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
	var outInfo *graph.ObjectInfo
	var err error
	var array bool
	if outInfo, array, err = checkOutArguments(svc, rMap, info.Response, log); err != nil {
		objtraver.DecReferenceCount(rMap, info.Response)
		return nil, err
	}
	// 返回数据的类型
	svc.ReturnType = svc.funcType.Out(0)
	svc.ReturnObjectID = outInfo.ID
	svc.ReturnObjectIsArray = array

	// 解析输入参数
	var inInfo *graph.ObjectInfo
	var inReq *request.RequestObjectField
	//var inObj *request.RequestObject
	inMap := map[string]int{}
	inInfo, inReq, err = checkInArguments(svc, requestObjectManager, inMap, info.Request, log)
	if err != nil {
		objtraver.DecReferenceCount(rMap, info.Response)
		objtraver.DecReferenceCount(inMap, info.Request)
		return nil, err
	}
	if inReq != nil {
		svc.requestObject = inReq
		svc.RequestObjectIsArray = inReq.Slice
	}
	if inInfo != nil {
		svc.RequestObjectID = inInfo.ID
	}
	return svc, nil
}

// 检查函数的输入参数
func checkInArguments(svc *JSONRunner,
	requestObjectManager *request.ObjectManager,
	referenceMap map[string]int,
	objMap map[string]*graph.ObjectInfo,
	log context.Log) (*graph.ObjectInfo, *request.RequestObjectField, error) {
	ic := svc.funcType.NumIn()
	var inInfo *graph.ObjectInfo
	var inObj *request.RequestObject
	var inObjField *request.RequestObjectField
	svc.inputArgs = make([]arguments.Argument, ic, ic)
	for n := 0; n < ic; n++ {
		in := svc.funcType.In(n)
		if types.IsRequirement(in) {
			// 用于判断必填字段检查的参数
			svc.inputArgs[n] = &arguments.ArgRequire{}
			continue
		}

		if types.IsResults(in) {
			svc.inputArgs[n] = &arguments.ArgResults{}
			continue
		}

		// 接口，必须是注入字段
		if in.Kind() == reflect.Interface {
			// 接口类型，必须是注入字段
			injector := svc.injectMgr.Find(in)
			if injector == nil {
				return nil, nil, fmt.Errorf("No injector exists with type %s: %s", in.Name(), svc.location)
			}
			svc.inputArgs[n] = &arguments.ArgInjector{
				Injector:    injector,
				IsInterface: true,
				ValueIsPtr:  true, // 接口类型，认为是指针类型
			}
			if injector.AccessController {
				// 权限控制相关的注入内容
				svc.AccessControllers = append(svc.AccessControllers, injector)
			}
			continue
		}

		typeInfo, err := reflects.ParseType(svc.location, in)
		if err != nil {
			return nil, nil, fmt.Errorf("Invalid service function : %s: %w", svc.location, err)
		}
		if n == 0 {
			// 结构体
			svc.loaderStruct = loader.ParseLoader(in, svc.injectMgr)
			svc.inputArgs[n] = svc.loaderStruct
			continue
		}

		if typeInfo.IsStruct {
			// 结构类型，必须定义为指针
			//if !typeInfo.IsRawPtr {
			//	return nil, nil, fmt.Errorf("A struct argument must be a pointer %s: %s", typeInfo.Name, svc.location)
			//}
			// 优先认定注入类型
			injector := svc.injectMgr.Find(typeInfo.Reference)
			if injector != nil {
				// 在注入类型中找到
				svc.inputArgs[n] = &arguments.ArgInjector{
					Injector:    injector,
					IsInterface: false,
					ValueIsPtr:  typeInfo.IsRawPtr,
				}
				if injector.AccessController {
					// 权限控制相关的注入内容
					svc.AccessControllers = append(svc.AccessControllers, injector)
				}
			} else {
				// 不是注入类型，那就一定是输入参数
				if svc.requestObject != nil {
					// 已经有输入参数了
					return nil, nil, fmt.Errorf("A service only has one request argument %s: %s", typeInfo.Name, svc.location)
				}

				var err error
				if inInfo, inObj, err = objtraver.Traversal(svc.location, typeInfo.Raw, referenceMap, objMap, requestObjectManager, log); err != nil {
					return inInfo, nil, fmt.Errorf("JSONRunner function invalid type: %s -> %s", err, svc.location)
				}

				inObjField = request.GenerateRequestObjectField(nil, "", typeInfo, false, func(err error) {
					log.Warn("%s: %s", svc.Name, err.Error())
				})

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
			if inInfo, inObj, err = objtraver.Traversal(svc.location, typeInfo.Raw, referenceMap, objMap, requestObjectManager, log); err != nil {
				return inInfo, nil, fmt.Errorf("JSONRunner function invalid type: %s -> %s", err, svc.location)
			}
			inObjField = request.GenerateRequestObjectField(nil, "", typeInfo, false, func(err error) {
				if err != nil {
					log.Warn("%s: %s", svc.Name, err.Error())
				}
			})
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
func checkOutArguments(svc *JSONRunner, referenceMap map[string]int, objMap map[string]*graph.ObjectInfo, log context.Log) (*graph.ObjectInfo, bool, error) {
	funcLoc := svc.location.String()
	oc := svc.funcType.NumOut()
	if oc != 2 {
		//log.Warn("JSONRunner function Must has return 2 values: %s", funcLoc)
		return nil, false, fmt.Errorf("JSONRunner function Must has return 2 values: %s", funcLoc)
	}

	o := svc.funcType.Out(1)
	if !types.IsError(o) {
		return nil, false, fmt.Errorf("JSONRunner function's second return argument must be type of error: %s", funcLoc)
	}
	//rMap := map[string]int{}
	o = svc.funcType.Out(0)
	outObj, _, err := objtraver.Traversal(svc.location, o, referenceMap, objMap, nil, log)
	if err != nil {
		//objtraver.DecReferenceCount(referenceMap, outObj)
		return nil, false, fmt.Errorf("JSONRunner function invalid type: %w -> %s", err, funcLoc)
	}
	// } else {
	//decReferenceCount(referenceMap, info)
	return outObj, o.Kind() == reflect.Slice, nil
	// }
}

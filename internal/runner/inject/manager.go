package inject

import (
	"fmt"
	"reflect"

	"github.com/seerx/runjson/internal/types"

	"github.com/seerx/runjson/internal/reflects"
)

// InjectorManager 注入函数管理
type InjectorManager struct {
	injectors map[reflect.Type]*Injector
}

// NewManager 创建注入管理器
func NewManager() *InjectorManager {
	return &InjectorManager{
		injectors: map[reflect.Type]*Injector{},
	}
}

// Find 查找注入函数
func (im *InjectorManager) Find(typ reflect.Type) *Injector {
	if inj, exists := im.injectors[typ]; exists {
		return inj
	}
	return nil
}

func (im *InjectorManager) Register(fn interface{}) error {
	return im.RegisterWithProxy(fn, nil, nil)
}

// RegisterWithProxy 注册注入函数
func (im *InjectorManager) RegisterWithProxy(fn interface{}, injectType reflect.Type, beenProxyFn interface{}) error {
	var loc *reflects.Location
	if beenProxyFn != nil && injectType != nil {
		loc = reflects.ParseFuncLocation(beenProxyFn)
	} else {
		loc = reflects.ParseFuncLocation(fn)
	}

	typ := reflect.TypeOf(fn)
	if typ.Kind() != reflect.Func {
		return fmt.Errorf("A injector must be a func [%s]", loc.String())
	}

	oc := typ.NumOut()
	if oc != 2 {
		// 返回参数必须是 2 个
		return fmt.Errorf("Injector func must return 2 args [%s]", loc.String())
	}
	if injectType == nil {
		// 如果没有指定 injectType ，使用注入函数的第一个返回值作为 injectType
		injectType = typ.Out(0)
	}

	// 判断 injectType 是否合乎规定
	injectTypeIsPtr := injectType.Kind() == reflect.Ptr
	if injectTypeIsPtr {
		// 是指针类型，获取实际类型
		injectType = injectType.Elem()
	}
	if injectType.Kind() != reflect.Interface && ((!injectTypeIsPtr) || injectType.Kind() != reflect.Struct) {
		// 如果注入类型不是接口，且不是指向结构的指针，则不支持注入
		return fmt.Errorf("Injector func first return value must be interface or a poniter of struct [%s]", loc.String())
	}

	// // 第一个返回值，必须是接口或者指向结构体的指针
	// if injectType == nil {
	// 	// ru如果没有指定类型
	// 	o1 := typ.Out(0)
	// 	o1Ptr := o1.Kind() == reflect.Ptr

	// 	injectType = o1
	// 	if o1Ptr {
	// 		injectType = o1.Elem()
	// 	}

	// 	if o1.Kind() != reflect.Interface && ((!o1Ptr) || injectType.Kind() != reflect.Struct) {
	// 		return fmt.Errorf("Injector func first return value must be interface or a poniter of struct [%s]", loc.String())
	// 	}
	// } else {
	// 	// 指定了注入类型
	// 	if injectType.Kind() == reflect.Ptr {
	// 		injectType = injectType.Elem()
	// 	}
	// }

	//o1Typ := reflects.ParseType(loc, o1)
	//if (o1.Kind() != reflect.Interface) && (o1.Kind() != reflect.Struct) {
	//	// 不是接口，且不是结构体
	//	return fmt.Errorf("Injector func first return value must be struct or interface [%s]", loc.String())
	//}
	// 第二个参数必须是 error
	o2 := typ.Out(1)
	if !types.IsError(o2) {
		// 不是 error
		return fmt.Errorf("Injector func second return must be error [%s]", loc.String())
	}

	// 查找是否已经存在
	old, exists := im.injectors[injectType]
	if exists {
		if old.Location.Equals(loc) {
			// 重复注册
			return nil
		}
		// 已经存在
		return fmt.Errorf("Inject type [%s] is Registered by func [%s]", injectType.Name(), old.Location.String())
	}

	// 判断输入参数
	ic := typ.NumIn()
	if ic != 1 {
		return fmt.Errorf("Injector func must recieve one argument [%s]", loc.String())
	}
	inType := typ.In(0)
	inTypeIsPtr := inType.Kind() == reflect.Ptr
	if inTypeIsPtr {
		inType = inType.Elem()
	}
	if !types.IsInjectParam(inType) {
		return fmt.Errorf("Injector func must recieve one argument of map[string]interface{} [%s]", loc.String())
	}

	// 注册
	im.injectors[injectType] = &Injector{
		Type:                  injectType,
		Func:                  reflect.ValueOf(fn),
		Location:              loc,
		ArgIsPtr:              inTypeIsPtr,
		ReturnTypeIsInterface: injectType.Kind() == reflect.Interface,
	}

	return nil
}

package inject

import (
	"fmt"
	"reflect"

	"github.com/seerx/chain/internal/types"

	"github.com/seerx/chain/internal/reflects"
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

// Register 注册注入函数
func (im *InjectorManager) Register(fn interface{}) error {
	loc := reflects.ParseFuncLocation(fn)

	typ := reflect.TypeOf(fn)
	if typ.Kind() != reflect.Func {
		return fmt.Errorf("A injector must be a func [%s]", loc.String())
	}

	oc := typ.NumOut()
	if oc != 2 {
		// 返回参数必须是 2 个
		return fmt.Errorf("Injector func must return 2 args [%s]", loc.String())
	}

	// 第一个返回值，必须是接口或者指向结构体的指针
	o1 := typ.Out(0)
	o1Typ := reflects.ParseType(loc, o1)
	if (!o1Typ.IsInterface) && (!(o1Typ.IsRawPtr && o1Typ.IsStruct)) {
		// 不是接口，且不是指向结构体的指针
		return fmt.Errorf("Injector func first return must be a pointer of struct or a interface [%s]", loc.String())
	}
	// 第二个参数必须是 error
	o2 := typ.Out(1)
	if !types.IsError(o2) {
		// 不是 error
		return fmt.Errorf("Injector func second return must be error [%s]", loc.String())
	}

	keyType := o1Typ.Reference

	// 查找是否已经存在
	old, exists := im.injectors[keyType]
	if exists {
		if old.Location.Equals(loc) {
			// 重复注册
			return nil
		}
		// 已经存在
		return fmt.Errorf("Type [%s] is Registered by func [%s]", o1Typ.Name, old.Location.String())
	}

	// 判断输入参数
	ic := typ.NumIn()
	if ic != 1 {
		return fmt.Errorf("Injector func must recieve one argument [%s]", loc.String())
	}
	inType := typ.In(0)
	if !types.IsInjectParam(inType) {
		return fmt.Errorf("Injector func must recieve one argument of map[string]interface{} [%s]", loc.String())
	}

	// 注册
	im.injectors[keyType] = &Injector{
		Type:     keyType,
		Func:     reflect.ValueOf(fn),
		Location: loc,
	}

	return nil
}

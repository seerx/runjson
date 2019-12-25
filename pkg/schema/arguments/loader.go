package arguments

import (
	"reflect"

	"github.com/seerx/runjson/pkg/inject"
)

var loaderMap = map[reflect.Type]*LoaderStruct{}

// LoaderStruct loader 结构定义
type LoaderStruct struct {
	Type reflect.Type
}

func (ls *LoaderStruct) CreateValue(ctx *ArgumentContext) reflect.Value {
	return reflect.New(ls.Type)
}

func (ls *LoaderStruct) IsInjectInterface() bool {
	return false
}

type LoaderInstance struct {
}

func (ls *LoaderStruct) New() *LoaderInstance {
	return &LoaderInstance{}
}

// ParseLoader 解析 Loader 结构
func ParseLoader(loaderType reflect.Type, injectorManager *inject.InjectorManager) *LoaderStruct {
	if ls, e := loaderMap[loaderType]; e {
		return ls
	}
	if loaderType.Kind() == reflect.Ptr {
		loaderType = loaderType.Elem()
	}
	ls := &LoaderStruct{Type: loaderType}
	loaderMap[loaderType] = ls
	return ls
}

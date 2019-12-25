package request

import "reflect"

// RequestObjectManager 输入对象池
type RequestObjectManager struct {
	objMap map[reflect.Type]*RequestObject
}

func NewRequestObjectManager() *RequestObjectManager {
	return &RequestObjectManager{
		objMap: map[reflect.Type]*RequestObject{},
	}
}

func (iop *RequestObjectManager) Find(typ reflect.Type) *RequestObject {
	if o, exists := iop.objMap[typ]; exists {
		return o
	}
	return nil
}

func (iop *RequestObjectManager) Register(in *RequestObject) {
	iop.objMap[in.Type] = in
}

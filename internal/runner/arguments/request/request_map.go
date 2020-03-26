package request

import "reflect"

// ObjectManager 输入对象池
type ObjectManager struct {
	objMap map[reflect.Type]*RequestObject
}

// NewRequestObjectManager 创建请求对象管理池
func NewRequestObjectManager() *ObjectManager {
	return &ObjectManager{
		objMap: map[reflect.Type]*RequestObject{},
	}
}

// Find 查找对象
func (iop *ObjectManager) Find(typ reflect.Type) *RequestObject {
	if o, exists := iop.objMap[typ]; exists {
		return o
	}
	return nil
}

// Register 注册对象
func (iop *ObjectManager) Register(in *RequestObject) {
	iop.objMap[in.Type] = in
}

package graph

// ObjectInfo 输入、输出数据信息描述
type ObjectInfo struct {
	ID          string        `json:"id,omitempty"`
	ReferenceID string        `json:"reference,omitempty"` // 指向其它 ObjectInfo.ID，原生数据类型为空
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Require     bool          `json:"require,omitempty"` // TODO 递归引用时，需要处理截断 Require 链
	Array       bool          `json:"array,omitempty"`   // 是否是数组
	Description string        `json:"description,omitempty"`
	Children    []*ObjectInfo `json:"children,omitempty"`
	Deprecated  bool          `json:"deprecated"` // 是否不建议使用

	ReferenceCount int `json:"-"` // 引用次数
}

package graph

// ServiceInfo 服务功能定义，给前端的接口
type ServiceInfo struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Deprecated     bool   `json:"deprecated"`
	InputObjectID  string `json:"inputObjectId,omitempty"`
	OutputObjectID string `json:"outputObjectId,omitempty"`
}

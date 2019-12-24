package apimap

// ServiceInfo 服务功能定义，给前端的接口
type ServiceInfo struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	InputObjectID  string `json:"inputObjectId"`
	OutputObjectID string `json:"outputObjectId"`
}

package chain

// Response 返回的数据类型
type Response struct {
	Error string `json:"error,omitempty"`
	Data  interface{}
}

type Responses map[string]*Response

// Request 请求对象
type Request struct {
	Service string      `json:"service"`
	Args    interface{} `json:"args"`
}

// 请求列表
type Requests []*Request
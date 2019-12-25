package runjson

// Response 返回的数据类型
type Response struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type Responses map[string]*Response

// Request 请求对象
type Request struct {
	Service string      `json:"service"`
	Alias   string      `json:"alias"` // 别名，用于接收返回数据
	Args    interface{} `json:"args"`
}

// 请求列表
type Requests []*Request

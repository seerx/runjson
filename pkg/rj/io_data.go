package rj

import "reflect"

// ResponseItem 返回项
type ResponseItem struct {
	Error    string       `json:"error,omitempty"`
	Data     interface{}  `json:"data,omitempty"`
	DataType reflect.Type `json:"-"`
}

// Response 返回值
type Response map[string][]*ResponseItem

// Request 请求对象
type Request struct {
	Service string      `json:"service"`
	Args    interface{} `json:"args"`
}

// Requests 请求列表
type Requests []*Request

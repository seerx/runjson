package rj

// ResponseContext 用于获取前面执行的结果
type ResponseContext interface {
	CallCount() int // 本次调用的 API 总数
	CallIndex() int // 当前 API 调用序列
	Get(method interface{}) ([]*ResponseItem, error)
}

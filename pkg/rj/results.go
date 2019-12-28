package rj

// Results 用于获取前面执行的结果
type Results interface {
	Get(method interface{}) ([]*ResponseItem, error)
}

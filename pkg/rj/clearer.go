package rj

// OnComplete 执行完成后，清理工作，用于注入对象
// 注入对象实现此接口，在任务执行完后会自动执行 Clear 函数
type OnComplete interface {
	Clear()
}

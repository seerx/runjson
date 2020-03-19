package rj

// GroupFunc 分组函数名称
const GroupFunc = "Group"

// Loader 承载接口
type Loader interface {
	Group() *Group
}

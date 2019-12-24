package intf

const GROUP_FUNC = "Group"

// Loader 承载接口
type Loader interface {
	Group() *Group
}

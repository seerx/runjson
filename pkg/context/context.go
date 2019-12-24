package context

import "context"

// Context 执行上下文
type Context struct {
	Context context.Context
	Param   map[string]interface{}
}

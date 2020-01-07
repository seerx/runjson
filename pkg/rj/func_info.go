package rj

import "github.com/seerx/runjson/pkg/graph"

// FuncInfo 功能信息
type FuncInfo struct {
	Description    string
	Deprecated     bool
	InputIsRequire bool
	History        []*graph.CR
}

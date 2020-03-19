package rj

// Require 必填接口
type Require interface {
	HasFields(fields ...string) bool
	Requires(fields ...string)
}

package intf

type Require interface {
	HasFields(fields ...string) bool
	Requires(fields ...string)
}

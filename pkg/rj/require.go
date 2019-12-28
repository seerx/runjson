package rj

type Require interface {
	HasFields(fields ...string) bool
	Requires(fields ...string)
}

package rj

// InjectArg 注入函数的参数
type InjectArg struct {
	Args     map[string]interface{}
	Service  string
	Response ResponseContext
}

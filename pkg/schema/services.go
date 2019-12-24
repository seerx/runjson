package schema

// Services 服务信息，用于执行服务
type Services struct {
	ServiceMap map[string]*Service
}

// AddService 添加服务
func (s *Services) AddService(svc *Service) {
	s.ServiceMap[svc.Name] = svc
}

func (s *Services) GetService(svcName string) *Service {
	return s.ServiceMap[svcName]
}

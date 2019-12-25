package graph

import (
	"fmt"
)

// GroupInfo 分组定义信息
type GroupInfo struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Services    []*ServiceInfo `json:"services"`
}

// GenerateServiceName 生成组内的服务名称
func (g *GroupInfo) GenerateServiceName(funcName string) (string, error) {
	name := fmt.Sprintf("%s.%s", g.Name, funcName)
	for _, svc := range g.Services {
		if svc.Name == name {
			// 存在相同名称的服务
			return name, fmt.Errorf("Service named %s is exists", name)
		}
	}
	return name, nil
}

func (g *GroupInfo) AddService(svc *ServiceInfo) {
	g.Services = append(g.Services, svc)
}

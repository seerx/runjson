package graph

// APIInfo 总体结构
type APIInfo struct {
	Groups   []*GroupInfo           `json:"groups"`
	Response map[string]*ObjectInfo `json:"response"`
	Request  map[string]*ObjectInfo `json:"request"`
}

// GetGroup 根据 rj.Group 获取 map 中的组信息
func (mi *APIInfo) GetGroup(grpName, grpInfo string) *GroupInfo {
	for _, group := range mi.Groups {
		if group.Name == grpName {
			// 组已经存在
			return group
		}
	}
	// 生成新的组信息
	group := &GroupInfo{
		Name:        grpName,
		Description: grpInfo,
	}
	mi.Groups = append(mi.Groups, group)
	return group
}

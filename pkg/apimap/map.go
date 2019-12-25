package apimap

import "github.com/seerx/runjson/pkg/intf"

// MapInfo 总体结构
type MapInfo struct {
	Groups   []*GroupInfo           `json:"groups"`
	Response map[string]*ObjectInfo `json:"response"`
	Request  map[string]*ObjectInfo `json:"request"`
}

// GetGroup 根据 intf.Group 获取 map 中的组信息
func (mi *MapInfo) GetGroup(grp *intf.Group) *GroupInfo {
	for _, group := range mi.Groups {
		if group.Name == grp.Name {
			// 组已经存在
			return group
		}
	}
	// 生成新的组信息
	group := &GroupInfo{
		Name:        grp.Name,
		Description: grp.Description,
	}
	mi.Groups = append(mi.Groups, group)
	return group
}

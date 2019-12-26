package tt

import (
	"github.com/seerx/runjson/pkg/intf"
)

type TT struct {
}

func (t TT) Group() *intf.Group {
	return &intf.Group{
		Name:        "QQ",
		Description: "123",
	}
}

func (t TT) New() (string, error) {
	return "123", nil
}

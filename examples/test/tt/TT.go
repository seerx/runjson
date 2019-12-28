package tt

import (
	"github.com/seerx/runjson/pkg/rj"
)

type TT struct {
}

func (t TT) Group() *rj.Group {
	return &rj.Group{
		Name:        "QQ",
		Description: "123",
	}
}

func (t TT) New() (string, error) {
	return "123", nil
}

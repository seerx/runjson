package chain

import (
	"github.com/seerx/chain/pkg/intf"
	"github.com/sirupsen/logrus"
)

var log = logrus.Logger{
	ReportCaller: true,
}

var chain = &Chain{
	log: log,
}

func Register(loader ...intf.Loader) {
	chain.Register(loader...)
}

func Explian() error {
	return chain.Explain()
}

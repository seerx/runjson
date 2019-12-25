package runjson

import (
	"github.com/seerx/runjson/pkg/intf"
	"github.com/sirupsen/logrus"
)

var log = logrus.Logger{
	ReportCaller: true,
}

var chain = &Runner{
	log: log,
}

func Register(loader ...intf.Loader) {
	chain.Register(loader...)
}

func Explian() error {
	return chain.Explain()
}

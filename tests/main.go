package main

import (
	"github.com/xuliangTang/athena/athena"
	"github.com/xuliangTang/athena/tests/classes"
	"github.com/xuliangTang/athena/tests/conf"
)

func main() {
	athena.Ignite().
		Load(athena.NewConfigModule(&conf.MyConf), athena.NewFuse(), athena.NewI18nModule()).
		Attach(athena.NewRateLimit()).
		Mount("v1", nil, classes.NewTestClass()).
		Launch()
}

package main

import (
	"github.com/xuliangTang/athena/athena"
	"github.com/xuliangTang/athena/tests/classes"
	"github.com/xuliangTang/athena/tests/conf"
	"github.com/xuliangTang/athena/tests/configurations"
)

func main() {
	athena.Ignite().
		Configuration(
			configurations.NewK8sMaps(),
			configurations.NewK8sHandler(),
			configurations.NewK8sConfig()).
		Load(athena.NewConfigModule(&conf.MyConf), athena.NewFuse(), athena.NewI18nModule()).
		Attach(athena.NewRateLimit()).
		Mount("v1", nil,
			classes.NewTestClass(),
			classes.NewK8sClass()).
		Launch()
}

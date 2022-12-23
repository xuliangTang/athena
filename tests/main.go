package main

import (
	"github.com/xuliangTang/athena/athena"
	"github.com/xuliangTang/athena/athena/middlewares"
	"github.com/xuliangTang/athena/athena/plugins"
	"github.com/xuliangTang/athena/tests/classes"
	"github.com/xuliangTang/athena/tests/configurations"
	"github.com/xuliangTang/athena/tests/properties"
)

func main() {
	athena.Ignite().
		Configuration(
			configurations.NewK8sMaps(),
			configurations.NewK8sHandler(),
			configurations.NewK8sConfig()).
		MappingConfig(&properties.MyConf).
		RegisterPlugin(plugins.NewI18n(), plugins.NewFuse()).
		Attach(middlewares.NewRateLimit()).
		Mount("v1", nil,
			classes.NewTestClass(),
			classes.NewK8sClass()).
		Launch()
}

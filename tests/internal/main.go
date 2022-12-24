package main

import (
	"github.com/xuliangTang/athena/athena"
	"github.com/xuliangTang/athena/athena/middlewares"
	"github.com/xuliangTang/athena/athena/plugins"
	classes2 "github.com/xuliangTang/athena/tests/internal/classes"
	configurations2 "github.com/xuliangTang/athena/tests/internal/configurations"
	"github.com/xuliangTang/athena/tests/internal/properties"
)

func main() {
	athena.Ignite().
		Configuration(
			configurations2.NewK8sMaps(),
			configurations2.NewK8sHandler(),
			configurations2.NewK8sConfig()).
		MappingConfig(&properties.MyConf).
		RegisterPlugin(plugins.NewI18n(), plugins.NewFuse()).
		Attach(middlewares.NewRateLimit()).
		Mount("v1", nil,
			classes2.NewTestClass(),
			classes2.NewK8sClass()).
		Launch()
}

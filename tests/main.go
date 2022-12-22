package main

import (
	"github.com/lain/athena/athena"
	"github.com/lain/athena/tests/classes"
)

func main() {
	athena.Ignite().
		Load(athena.NewFuse(), athena.NewI18nModule()).
		Attach(athena.NewRateLimit()).
		Mount("v1", nil, classes.NewTestClass()).
		Launch()
}

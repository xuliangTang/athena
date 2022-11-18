package main

import (
	"github.com/XNXKTech/athena/athena"
	"github.com/XNXKTech/athena/tests/classes"
)

func main() {
	athena.Ignite().
		Load(athena.NewFuse()).
		Attach(athena.NewRateLimit()).
		Mount("v1", nil, classes.NewTestClass()).
		Launch()
}

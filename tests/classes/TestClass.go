package classes

import (
	"github.com/XNXKTech/athena/athena"
	"github.com/gin-gonic/gin"
)

type TestClass struct {
}

func NewTestClass() *TestClass {
	return &TestClass{}
}

func (this *TestClass) test(ctx *gin.Context) *athena.Json {
	return &athena.Json{
		"message": "test",
	}
}

func (this *TestClass) Build(athena *athena.Athena) {
	athena.Handle("GET", "/test", this.test)
}

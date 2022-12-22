package classes

import (
	"github.com/lain/athena/athena"
	"github.com/gin-gonic/gin"
)
{{$ClassName:=(printf "%s%s" .ControllerName "Class") | Ucfirst}}
type {{$ClassName}} struct {
}

func New{{$ClassName}}() *{{$ClassName}} {
	return &{{$ClassName}}{}
}

func(this *{{$ClassName}}) {{.ControllerName}}Detail(ctx *gin.Context) *athena.Json {
	return &athena.Json{}
}

func(this *{{$ClassName}}) Build(athena *athena.Athena){
	// athena.Handle("GET","/your path/:id",this.{{.ControllerName}}Detail)
}

package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuliangTang/athena/athena/lib"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
)

// ErrorCache @Middleware 错误捕获
type ErrorCache struct{}

func NewErrorCache() *ErrorCache {
	return &ErrorCache{}
}

func (*ErrorCache) OnRequest(ctx *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			var errInfo string
			switch e := e.(type) {
			case string:
				errInfo = e
			case runtime.Error:
				errInfo = e.Error()
			case error:
				errInfo = e.Error()
			default:
				errInfo = fmt.Sprintf("unknown error type: %s", reflect.TypeOf(e).String())
			}

			printError(errInfo)
			lib.Logger().Error("panic: " + errInfo)

			code := http.StatusBadRequest
			if ctxCode, exist := ctx.Get("athena_httpStatusCode"); exist {
				if v, ok := ctxCode.(int); ok {
					code = v
				}
			}

			ctx.AbortWithStatusJSON(code, gin.H{"error": e})
		}
	}()

	ctx.Next()
}

func printError(err interface{}) {
	if os.Getenv("GIN_MODE") == "release" {
		return
	}
	log.Println(err)
	log.Println(lib.Logger().GetStack(10))
}

package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuliangTang/athena/athena/lib"
	"net/http"
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

			lib.Logger().Error("panic: " + errInfo)

			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": e})
		}
	}()

	ctx.Next()
}

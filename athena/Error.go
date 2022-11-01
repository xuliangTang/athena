package athena

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"runtime"
)

func ErrorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
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

				Logger().Error("PANIC ERROR", zap.String("info", errInfo))
				context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": e})
			}
		}()
		context.Next()
	}
}

func Error(err error, msg ...string) {
	if err != nil {
		errMsg := err.Error()
		if len(msg) > 0 {
			errMsg = msg[0]
		}
		panic(errMsg)
	}
}

func Unwrap(result any, err error) any {
	if err != nil {
		panic(err.Error())
	}

	return result
}

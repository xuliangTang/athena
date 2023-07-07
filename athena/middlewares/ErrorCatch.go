package middlewares

import (
	"cuelang.org/go/cue"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuliangTang/athena/athena/lib"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
)

// ErrorCatch @Middleware 错误捕获
type ErrorCatch struct {
	RspErrorValue cue.Value
}

func NewErrorCatch(rspErrorValue cue.Value) *ErrorCatch {
	return &ErrorCatch{RspErrorValue: rspErrorValue}
}

const (
	cueOutput = "output"
)

func (this *ErrorCatch) OnRequest(ctx *gin.Context) {
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

			this.printError(errInfo)
			lib.Logger().Error("panic: " + errInfo)

			code := http.StatusBadRequest
			if ctxCode, exist := ctx.Get("athena_httpStatusCode"); exist {
				if v, ok := ctxCode.(int); ok {
					code = v
				}
			}

			cueValue := this.genRspByCue(code, errInfo) // 生成响应内容
			ctx.AbortWithStatusJSON(code, cueValue)
		}
	}()

	ctx.Next()
}

func (this *ErrorCatch) printError(err interface{}) {
	if os.Getenv("GIN_MODE") == "release" {
		return
	}
	log.Println(err)
	log.Println(lib.Logger().GetStack(10))
}

// 根据cue模板生成响应结构体
func (this *ErrorCatch) genRspByCue(code int, message string) cue.Value {
	errV := this.RspErrorValue
	
	// 遍历模板节点
	if field, err := errV.LookupPath(cue.ParsePath(cueOutput)).Fields(); err == nil {
		for field.Next() {
			// 获取每个节点的注释
			parsePath := cue.ParsePath(cueOutput + "." + field.Label())
			doc := errV.LookupPath(parsePath).Doc()

			// 替换
			for _, d := range doc {
				getDoc := strings.TrimSpace(d.Text())

				if getDoc == "@code" {
					errV = errV.FillPath(parsePath, code)
					break
				} else if getDoc == "@message" {
					errV = errV.FillPath(parsePath, message)
					break
				}
			}
		}
	}

	return errV.LookupPath(cue.ParsePath(cueOutput)).Value()
}

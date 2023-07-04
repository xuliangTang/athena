package middlewares

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuliangTang/athena/athena/config"
	"github.com/xuliangTang/athena/athena/lib"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
)

// ErrorCatch @Middleware 错误捕获
type ErrorCatch struct{}

func NewErrorCatch() *ErrorCatch {
	return &ErrorCatch{}
}

const (
	cueOutput = "output"
)

func (*ErrorCatch) OnRequest(ctx *gin.Context) {
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

			cueValue := genRspByCue(code, errInfo) // 生成响应内容
			ctx.AbortWithStatusJSON(code, cueValue)
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

// 根据cue模板生成响应结构体
func genRspByCue(code int, message string) cue.Value {
	// 生成模板的value
	cc := cuecontext.New()
	v := cc.CompileString(config.AppConf.RspCueTpl.ErrorTpl)

	if field, err := v.LookupPath(cue.ParsePath(cueOutput)).Fields(); err == nil { // 遍历模板节点
		for field.Next() {
			// 获取每个节点的注释
			parsePath := cue.ParsePath(cueOutput + "." + field.Label())
			doc := v.LookupPath(parsePath).Doc()

			// 替换
			for _, d := range doc {
				getDoc := strings.TrimSpace(d.Text())

				if getDoc == "@code" {
					v = v.FillPath(parsePath, code)
					break
				} else if getDoc == "@message" {
					v = v.FillPath(parsePath, message)
					break
				}
			}
		}
	}

	return v.LookupPath(cue.ParsePath(cueOutput)).Value()
}

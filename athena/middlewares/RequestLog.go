package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuliangTang/athena/athena/lib"
	"go.uber.org/zap"
	"time"
)

// RequestLog @Middleware 记录请求日志
type RequestLog struct{}

func NewRequestLog() *RequestLog {
	return &RequestLog{}
}

func (*RequestLog) OnRequest(ctx *gin.Context) {
	startTime := time.Now()
	ctx.Next()
	endTime := time.Now()
	execTime := endTime.Sub(startTime) // 响应时间
	requestMethod := ctx.Request.Method
	requestURI := ctx.Request.RequestURI
	statusCode := ctx.Writer.Status()
	requestIP := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()
	lib.Logger().Info(fmt.Sprintf("%s - %s %s[%d]", requestIP, requestMethod, requestURI, statusCode),
		zap.String("execTime", execTime.String()),
		zap.String("userAgent", userAgent),
	)

	ctx.Next()
}

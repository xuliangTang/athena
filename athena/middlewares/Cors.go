package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors @Middleware 跨域中间件
type Cors struct{}

func NewCors() *Cors {
	return &Cors{}
}

func (*Cors) OnRequest(ctx *gin.Context) {
	method := ctx.Request.Method
	origin := ctx.Request.Header.Get("Origin")
	if origin != "" {
		ctx.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, X-Token")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")
	}
	if method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
	}

	ctx.Next()
}

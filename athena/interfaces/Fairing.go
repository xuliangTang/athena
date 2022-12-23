package interfaces

import "github.com/gin-gonic/gin"

type IFairing interface {
	OnRequest(*gin.Context)
}

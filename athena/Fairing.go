package athena

import "github.com/gin-gonic/gin"

type IFairing interface {
	OnRequest(*gin.Context) error
}

type ILoad interface {
	Run() error
}

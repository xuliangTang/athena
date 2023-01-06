package athena

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"sync"
)

const (
	CtxHttpStatusCode = "athena_httpStatusCode"
	CtxCode           = "athena_code"
	CtxMessage        = "athena_message"
)

var ResponderList []Responder
var ResponsePool *sync.Pool

type Responder interface {
	RespondTo() gin.HandlerFunc
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func init() {
	ResponderList = []Responder{
		new(StringResponder),
		new(AnyResponder),
		new(JsonResponder),
		new(CollectionResponder),
		new(VoidResponder),
	}

	ResponsePool = &sync.Pool{New: func() any {
		return &Response{
			Code:    http.StatusOK,
			Message: "success",
			Data:    nil,
		}
	}}
}

func Convert(handler interface{}) gin.HandlerFunc {
	handlerRef := reflect.ValueOf(handler)
	for _, responder := range ResponderList {
		responderRef := reflect.ValueOf(responder).Elem()
		if handlerRef.Type().ConvertibleTo(responderRef.Type()) {
			responderRef.Set(handlerRef)
			return responderRef.Interface().(Responder).RespondTo()
		}
	}
	return nil
}

// 从response池中拿出一个对象
func getResponse() *Response {
	return ResponsePool.Get().(*Response)
}

// 放回response对象池
func putResponse(response *Response) {
	ResponsePool.Put(response)
}

func getCode(ctx *gin.Context) (code int) {
	if ctxCode, exist := ctx.Get(CtxCode); exist {
		if v, ok := ctxCode.(int); ok {
			code = v
		}
	}

	return
}

func getHttpStatusCode(ctx *gin.Context) (code int) {
	code = http.StatusOK
	if ctxCode, exist := ctx.Get(CtxHttpStatusCode); exist {
		if v, ok := ctxCode.(int); ok {
			code = v
		}
	}

	return
}

func getMessage(ctx *gin.Context) (msg string) {
	msg = "success"
	if ctxCode, exist := ctx.Get(CtxMessage); exist {
		if v, ok := ctxCode.(string); ok {
			msg = v
		}
	}

	return
}

// Controller return any

type AnyResponder func(ctx *gin.Context) any

func (this AnyResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		response := getResponse()
		defer putResponse(response)
		response.Data = this(context)
		response.Message = getMessage(context)
		response.Code = getCode(context)
		context.JSON(getHttpStatusCode(context), response)
	}
}

// Controller return Json

type Json map[string]any
type JsonResponder func(*gin.Context) Json

func (this JsonResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		response := getResponse()
		defer putResponse(response)
		response.Data = this(context)
		response.Message = getMessage(context)
		response.Code = getCode(context)
		context.JSON(getHttpStatusCode(context), response)
	}
}

// Controller return Collection

type CollectionResponder func(ctx *gin.Context) Collection

func (this CollectionResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		response := getResponse()
		defer putResponse(response)
		response.Data = this(context)
		response.Message = getMessage(context)
		response.Code = getCode(context)
		context.JSON(getHttpStatusCode(context), response)
	}
}

// Controller return void

type Void struct{}
type VoidResponder func(ctx *gin.Context) Void

func (this VoidResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		this(context)
		context.Status(getHttpStatusCode(context))
	}
}

// Controller return string

type StringResponder func(*gin.Context) string

func (this StringResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		msg := this(context)
		context.String(getHttpStatusCode(context), msg)
	}
}

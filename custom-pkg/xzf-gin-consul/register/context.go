package register

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// Meta 元数据维护，用于在gin之间传递数据，如caller,code,message等

var metaKey = "__meta__"

type Code interface {
	GetCode() int
}

type Context struct {
	*gin.Context
}

func NewContext(ctx *gin.Context) *Context {
	return &Context{ctx}
}

func (ctx *Context) setMeta(key string, value interface{}) {
	if ctx == nil {
		return
	}
	meta := ctx.GetStringMap(metaKey)
	if meta == nil {
		meta = map[string]interface{}{}
		ctx.Set(metaKey, meta)
	}
	meta[key] = value
}

func (ctx *Context) getMeta(key string) interface{} {
	if ctx == nil {
		return nil
	}
	meta := ctx.GetStringMap(metaKey)
	if meta == nil {
		return nil
	}
	return meta[key]
}

func (ctx *Context) SetCaller(caller string) {
	ctx.setMeta("caller", caller)
}

func (ctx *Context) GetCaller() string {
	caller, _ := ctx.getMeta("caller").(string)
	return caller
}

func (ctx *Context) setCode(code string) {
	ctx.setMeta("code", code)
}

func (ctx *Context) GetMessage() string {
	message, _ := ctx.getMeta("message").(string)
	return message
}

func (ctx *Context) JSON(code int, obj interface{}) {
	if o, ok := obj.(Code); ok {
		ctx.setCode(strconv.Itoa(o.GetCode())) // auto set code into meta
	}
	ctx.Context.JSON(code, obj)
}

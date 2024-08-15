package register

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

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

func (ctx *Context) SetMeta(key string, value interface{}) {
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

func (ctx *Context) GetMeta(key string) interface{} {
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
	ctx.SetMeta("caller", caller)
}

func (ctx *Context) GetCaller() string {
	caller, _ := ctx.GetMeta("caller").(string)
	return caller
}

func (ctx *Context) SetCode(code string) {
	ctx.SetMeta("code", code)
}

func (ctx *Context) GetMessage() string {
	message, _ := ctx.GetMeta("message").(string)
	return message
}

func (ctx *Context) JSON(code int, obj interface{}) {
	if o, ok := obj.(Code); ok {
		ctx.SetCode(strconv.Itoa(o.GetCode())) // auto set code into meta
	}
	ctx.Context.JSON(code, obj)
}

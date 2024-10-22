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

func (ctx *Context) SetCaller(caller string) *Context {
	ctx.setMeta("caller", caller)
	return ctx
}

func (ctx *Context) GetCaller() string {
	caller, _ := ctx.getMeta("caller").(string)
	return caller
}

func (ctx *Context) SetError(error string) *Context {
	ctx.setMeta("error", error)
	return ctx
}
func (ctx *Context) GetError() string {
	err, _ := ctx.getMeta("error").(string)
	return err
}

func (ctx *Context) SetEmail(email string) *Context {
	ctx.setMeta("email", email)
	return ctx
}
func (ctx *Context) GetEmail() string {
	email, _ := ctx.getMeta("email").(string)
	return email
}

func (ctx *Context) SetMethod(method string) *Context {
	ctx.setMeta("method", method)
	return ctx
}
func (ctx *Context) GetMethod() string {
	method, _ := ctx.getMeta("method").(string)
	return method
}

func (ctx *Context) SetPassword(password string) *Context {
	ctx.setMeta("password", password)
	return ctx
}
func (ctx *Context) GetPassword() string {
	password, _ := ctx.getMeta("password").(string)
	return password
}

func (ctx *Context) SetUid(uid string) *Context {
	ctx.setMeta("uid", uid)
	return ctx
}
func (ctx *Context) GetUid() string {
	uid, _ := ctx.getMeta("uid").(string)
	return uid
}

func (ctx *Context) SetPage(page string) *Context {
	ctx.setMeta("page", page)
	return ctx
}
func (ctx *Context) GetPage() string {
	page, _ := ctx.getMeta("page").(string)
	return page
}

func (ctx *Context) SetPageSize(pageSize string) *Context {
	ctx.setMeta("pageSize", pageSize)
	return ctx
}
func (ctx *Context) GetPageSize() string {
	pageSize, _ := ctx.getMeta("pageSize").(string)
	return pageSize
}

func (ctx *Context) SetMuid(muid string) *Context {
	ctx.setMeta("muid", muid)
	return ctx
}
func (ctx *Context) GetMuid() string {
	muid, _ := ctx.getMeta("muid").(string)
	return muid
}

func (ctx *Context) SetSuid(suid string) *Context {
	ctx.setMeta("suid", suid)
	return ctx
}
func (ctx *Context) GetSuid() string {
	suid, _ := ctx.getMeta("suid").(string)
	return suid
}

func (ctx *Context) SetVisibility(visibility string) *Context {
	ctx.setMeta("visibility", visibility)
	return ctx
}
func (ctx *Context) GetVisibility() string {
	visibility, _ := ctx.getMeta("visibility").(string)
	return visibility
}

func (ctx *Context) SetTyp(typ string) *Context {
	ctx.setMeta("typ", typ)
	return ctx
}
func (ctx *Context) GetTyp() string {
	typ, _ := ctx.getMeta("typ").(string)
	return typ
}

func (ctx *Context) SetAid(aid string) *Context {
	ctx.setMeta("aid", aid)
	return ctx
}
func (ctx *Context) GetAid() string {
	aid, _ := ctx.getMeta("aid").(string)
	return aid
}

// auto set code into meta
func (ctx *Context) setCode(code string) {
	ctx.setMeta("code", code)
}

func (ctx *Context) GetCode() string {
	code, _ := ctx.getMeta("code").(string)
	return code
}

func (ctx *Context) JSON(code int, obj interface{}) {
	if o, ok := obj.(Code); ok {
		ctx.setCode(strconv.Itoa(o.GetCode())) // auto set code into meta
	}
	ctx.Context.JSON(code, obj)
}

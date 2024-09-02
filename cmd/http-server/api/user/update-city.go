package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
)

// updateProvinceCityParams 面向接口
type updateProvinceCityParams struct {
	Uid      string `json:"uid" form:"uid" binding:"required"`
	Province string `json:"province" form:"province" binding:"required"`
	City     string `json:"city" form:"city" binding:"required"`
}

// UpdateProvinceCity 更新用户的户籍省市
// 4000260
// 2000260
func (s *Server) UpdateProvinceCity(ctx *gin.Context) {
	p := &updateProvinceCityParams{}
	if err := ctx.Bind(p); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000260, Message: "参数错误"})
		return
	}
	u, err := user.GetUser(p.Uid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000261, Message: "获取用户信息失败"})
		return
	}
	pc := p.Province + " / " + p.City
	if pc == u.ProvinceCity {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000262, Message: "省市未改变"})
		return
	}
	if err = user.UpdateProvinceCity(p.Uid, pc); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000263, Message: "更新省市失败"})
		return
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000260, Message: "更新省市成功"})
	return
}

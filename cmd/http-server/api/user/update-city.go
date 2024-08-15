package user

import (
	"github.com/gin-gonic/gin"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
)

// updateProvinceCityParams 面向接口
type updateProvinceCityParams struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Province string `json:"province" form:"province" binding:"required"`
	City     string `json:"city" form:"city" binding:"required"`
}

// UpdateProvinceCity 更新用户的户籍省市
// 4000260
// 2000260
func (s *Server) UpdateProvinceCity(ctx *gin.Context) {
	p := &updateProvinceCityParams{}
	if err := ctx.Bind(p); err != nil {
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000260, Message: "bind params err"})
		return
	}
	u, err := user.GetUser(p.Email)
	if err != nil {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000261, Message: "get user err"})
		return
	}
	pc := p.Province + " / " + p.City
	if pc == u.ProvinceCity {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000262, Message: "province city is same"})
		return
	}
	if err := user.UpdateProvinceCity(p.Email, pc); err != nil {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000263, Message: "update province city err"})
		return
	}
	ctx.JSON(200, ginConsulRegister.Response{Code: 2000260, Message: "update province city succ"})
	return
}

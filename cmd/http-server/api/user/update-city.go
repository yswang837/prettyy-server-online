package user

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
	"prettyy-server-online/utils/metrics"
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
func (s *Server) UpdateProvinceCity(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("update-province-city", "total")
	p := &updateProvinceCityParams{}
	if err := ctx.Bind(p); err != nil {
		metrics.CommonCounter.Inc("update-province-city", "params-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000260, Message: "参数错误"})
		return
	}
	pc := p.Province + " / " + p.City
	ctx.SetUid(p.Uid).SetProvinceCity(pc)
	u, err := user.GetUser(p.Uid)
	if err != nil {
		metrics.CommonCounter.Inc("update-province-city", "get-user-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000261, Message: "获取用户信息失败"})
		return
	}
	if pc == u.ProvinceCity {
		metrics.CommonCounter.Inc("update-province-city", "same-province-city")
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000262, Message: "省市未改变"})
		return
	}
	if err = user.UpdateProvinceCity(p.Uid, pc); err != nil {
		metrics.CommonCounter.Inc("update-province-city", "update-province-city-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000263, Message: "更新省市失败"})
		return
	}
	metrics.CommonCounter.Inc("update-province-city", "succ")
	ctx.JSON(http.StatusOK, &ginConsulRegister.Response{Code: 2000260, Message: "更新省市成功"})
	return
}

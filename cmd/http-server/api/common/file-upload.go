package common

import (
	"github.com/gin-gonic/gin"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	xzf_qiniu "prettyy-server-online/custom-pkg/xzf-qiniu"
)

// FileUpload 文件上传
// 4000140
// 2000140

func (s *Server) FileUpload(ctx *gin.Context) {
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000140, Message: "参数绑定错误"})
		return
	}
	fileSize := fileHeader.Size
	url, err := xzf_qiniu.UploadFile(file, fileSize)
	if err != nil {
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000141, Message: "上传文件错误" + err.Error()})
		return
	}
	ctx.JSON(200, ginConsulRegister.Response{Code: 2000140, Message: "上传文件成功", Result: url})
	return
}
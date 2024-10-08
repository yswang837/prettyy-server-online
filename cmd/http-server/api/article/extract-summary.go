package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
)

// 提取文章摘要信息
// 2000380
// 4000380

const (
	qwenURL = "http://10.2.8.116:11434/api/generate"
)

type extractSummaryParams struct {
	Content string `json:"content" form:"content" binding:"required"` // 文章内容，用于提取摘要
}

type extractSummaryResp struct {
	Model              string  `json:"model,omitempty"`
	CreatedAt          string  `json:"created_at,omitempty"`
	Response           string  `json:"response,omitempty"`
	Done               bool    `json:"done,omitempty"`
	DoneReason         string  `json:"done_reason,omitempty"`
	Context            []int64 `json:"context,omitempty"`
	TotalDuration      int64   `json:"total_duration,omitempty"`
	LoadDuration       int64   `json:"load_duration,omitempty"`
	PromptEvalCount    int64   `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64   `json:"prompt_eval_duration,omitempty"`
	EvalCount          int64   `json:"eval_count,omitempty"`
	EvalDuration       int64   `json:"eval_duration,omitempty"`
}

func (s *Server) ExtractSummary(ctx *gin.Context) {
	params := &extractSummaryParams{}
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000380, Message: "参数错误"})
		return
	}
	query := map[string]interface{}{
		"model":  "qwen2:0.5b",
		"prompt": params.Content,
		"stream": false,
	}
	resp := &extractSummaryResp{}
	_, err := s.client.R().SetResult(resp).SetBody(query).Post(qwenURL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000381, Message: "调用通义千问出错"})
		return
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000380, Message: "调用qwen2:0.5b成功", Result: resp.Response})
	return
}

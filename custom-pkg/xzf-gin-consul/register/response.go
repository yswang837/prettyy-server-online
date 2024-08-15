package register

// Response is a default response format
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result,omitempty"`
}

func (r *Response) GetCode() int {
	if r == nil {
		return 0
	}
	return r.Code
}

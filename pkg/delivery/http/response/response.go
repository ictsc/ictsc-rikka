package response

import "github.com/gin-gonic/gin"

type response struct {
	Code    int         `json:"code"`
	Message *string     `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func JSON(ctx *gin.Context, code int, message string, data, err interface{}) {
	res := response{
		Code:  code,
		Data:  data,
		Error: err,
	}

	if message != "" {
		res.Message = &message
	}

	ctx.JSON(code, res)
}

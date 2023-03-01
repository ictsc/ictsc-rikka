package response

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type response struct {
	Code    int         `json:"code"`
	Message *string     `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type responseRaw struct {
	Code    int             `json:"code"`
	Message *string         `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   interface{}     `json:"error,omitempty"`
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

func JSONRaw(ctx *gin.Context, code int, message string, jsonData []byte, err interface{}) {
	res := responseRaw{
		Code:  code,
		Data:  jsonData,
		Error: err,
	}
	if message != "" {
		res.Message = &message
	}

	ctx.JSON(code, res)
}

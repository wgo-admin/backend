package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/errno"
)

// 定义发生错误时的返回消息体
type ErrResponse struct {
	Code    string `json:"code"`    // 业务错误码
	Message string `json:"message"` // 对外展示错误信息
}

// 将错误或响应数据写入 HTTP 响应体
// 使用 errno.Decode 方法，将我们业务层的错误信息进行解析并响应
func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		statusCode, _, message := errno.Decode(err)
		c.JSON(statusCode, responseBody{
			Code:    -1,
			Message: message,
			Success: false,
		})
		return
	}

	c.JSON(http.StatusOK, responseBody{
		Code:    0,
		Message: "本次请求成功",
		Success: true,
		Data:    data,
	})
}

type responseBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseOk(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, responseBody{
		Code:    0,
		Message: msg,
		Success: true,
		Data:    data,
	})
}

func ResponseFail(c *gin.Context, err error) {
	_, _, message := errno.Decode(err)
	c.JSON(http.StatusOK, responseBody{
		Code:    -1,
		Message: message,
		Success: false,
	})
}

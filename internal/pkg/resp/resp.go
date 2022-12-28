package resp

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
func Write(c *gin.Context, err error, data interface{}) {
	if err != nil {
		statusCode, errCode, message := errno.Decode(err)
		c.JSON(statusCode, ErrResponse{
			Code:    errCode,
			Message: message,
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
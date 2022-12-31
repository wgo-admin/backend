package user

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/known"
	"github.com/wgo-admin/backend/internal/pkg/log"
)

// login 登录控制层
func (ctrl *UserController) logout(c *gin.Context) {
	log.C(c.Request.Context()).Infow("Logout function called")

	username := c.GetString(known.XUsernameKey)

	// 调用业务层的 Login 方法
	err := ctrl.biz.User().Logout(c.Request.Context(), username)
	if err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "登出成功", nil)
}

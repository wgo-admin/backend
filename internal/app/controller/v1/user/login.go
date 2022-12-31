package user

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/log"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
	"github.com/wgo-admin/backend/pkg/validate"
)

// login 登录控制层
func (ctrl *UserController) login(c *gin.Context) {
	log.C(c.Request.Context()).Infow("Login function called")

	// 绑定参数
	var req v1.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.ResponseFail(c, errno.ErrBind)
		return
	}

	// 检验参数
	if err := validate.ValidateStruct(&req); err != nil {
		core.ResponseFail(c, errno.ErrInvalidParameter)
		log.C(c.Request.Context()).Errorw("invalid parameter", "err", err.Error())
		return
	}

	// 调用业务层的 Login 方法
	resp, err := ctrl.biz.User().Login(c.Request.Context(), &req)
	if err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "登录成功", resp)
}

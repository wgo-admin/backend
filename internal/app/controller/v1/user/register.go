package user

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/log"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
	"github.com/wgo-admin/backend/pkg/validate"
)

// 注册控制层方法
func (ctrl *UserController) register(c *gin.Context) {
	log.C(c.Request.Context()).Infow("Create user function called")

	var req v1.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.ResponseFail(c, errno.ErrBind)
		return
	}

	// 检验参数
	if err := validate.ValidateStruct(req); err != nil {
		core.ResponseFail(c, errno.ErrInvalidParameter)
		log.C(c.Request.Context()).Errorw("invalid parameter", "err", err.Error())
		return
	}

	if err := ctrl.biz.User().Register(c.Request.Context(), &req); err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "注册成功", nil)
}

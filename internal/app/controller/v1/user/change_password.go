package user

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/log"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
	"github.com/wgo-admin/backend/pkg/validate"
)

func (ctrl *UserController) changePassword(c *gin.Context) {
	log.C(c).Infow("ChangePassword user function called")

	var req v1.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.ResponseFail(c, errno.ErrBind)
		return
	}

	if err := validate.ValidateStruct(&req); err != nil {
		core.ResponseFail(c, errno.ErrInvalidParameter)
		return
	}

	if err := ctrl.biz.User().ChangePassword(c, c.Param("username"), &req); err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "修改成功", nil)
}

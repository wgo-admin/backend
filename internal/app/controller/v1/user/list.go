package user

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/log"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
)

func (ctrl *UserController) list(c *gin.Context) {
	log.C(c).Infow("List user function called")

	var req v1.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.ResponseFail(c, errno.ErrBind)
		return
	}

	resp, err := ctrl.biz.User().List(c, &req)
	if err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "查询成功", resp)
}

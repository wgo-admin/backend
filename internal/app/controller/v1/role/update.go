package role

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/log"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
	"github.com/wgo-admin/backend/pkg/validate"
)

func (ctrl *RoleController) update(c *gin.Context) {
	log.C(c.Request.Context()).Infow("update role function called")

	var req v1.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.ResponseFail(c, errno.ErrBind)
		return
	}

	if err := validate.ValidateStruct(req); err != nil {
		core.ResponseFail(c, errno.ErrInvalidParameter)
		return
	}

	if err := ctrl.biz.Role().Update(c.Request.Context(), c.Param("id"), &req); err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "更新成功", nil)
}

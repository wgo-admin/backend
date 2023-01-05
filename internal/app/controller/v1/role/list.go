package role

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/log"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
	"github.com/wgo-admin/backend/pkg/validate"
)

// 查询 role 列表
func (ctrl *RoleController) list(c *gin.Context) {
	log.C(c.Request.Context()).Infow("List role function called")

	var req v1.QueryRoleRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.ResponseFail(c, errno.ErrBind)
		return
	}

	if err := validate.ValidateStruct(req); err != nil {
		core.ResponseFail(c, errno.ErrInvalidParameter)
		return
	}

	resp, err := ctrl.biz.Role().List(c.Request.Context(), &req)
	if err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "查询成功", resp)
}

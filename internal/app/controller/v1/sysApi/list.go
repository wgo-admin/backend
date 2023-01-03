package sysApi

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/log"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
	"github.com/wgo-admin/backend/pkg/validate"
)

// 查询 sysApi 列表
func (ctrl *SysApiController) list(c *gin.Context) {
	log.C(c).Infow("List sysApi function called")

	var req v1.QuerySysApiRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.ResponseFail(c, errno.ErrBind)
		return
	}

	if err := validate.ValidateStruct(req); err != nil {
		core.ResponseFail(c, errno.ErrInvalidParameter)
		return
	}

	resp, err := ctrl.biz.SysApi().List(c, &req)
	if err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "查询成功", resp)
}

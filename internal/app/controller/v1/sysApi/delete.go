package sysApi

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/log"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
)

// 批量删除 sysApi
func (ctrl *SysApiController) delete(c *gin.Context) {
	log.C(c).Infow("delete sysApi function called")

	var req v1.IdsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.ResponseFail(c, errno.ErrBind)
		return
	}

	if err := ctrl.biz.SysApi().BatchDelete(c, &req); err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "删除成功", nil)
}

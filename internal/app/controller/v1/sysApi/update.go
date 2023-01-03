package sysApi

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/log"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
	"github.com/wgo-admin/backend/pkg/validate"
)

func (ctrl *SysApiController) update(c *gin.Context) {
	log.C(c).Infow("update sysApi function called")

	var req v1.UpdateSysApiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.ResponseFail(c, errno.ErrBind)
		return
	}

	if err := validate.ValidateStruct(req); err != nil {
		core.ResponseFail(c, errno.ErrInvalidParameter)
		return
	}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := ctrl.biz.SysApi().Update(c, id, &req); err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "更新成功", nil)
}

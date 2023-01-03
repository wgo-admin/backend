package sysApi

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
)

func (ctrl *SysApiController) get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	sysApi, err := ctrl.biz.SysApi().Get(c, id)
	if err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "查询成功", sysApi)
}

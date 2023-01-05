package menu

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/log"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
	"github.com/wgo-admin/backend/pkg/validate"
)

func (ctrl *MenuController) update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req v1.CreateOrUpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.ResponseFail(c, errno.ErrBind)
		return
	}

	if err := validate.ValidateStruct(req); err != nil {
		log.C(c).Errorw("ErrInvalidParameter", "error", err)
		core.ResponseFail(c, errno.ErrInvalidParameter)
		return
	}

	if err := ctrl.biz.Menu().Update(c, id, &req); err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "更新成功", nil)
}

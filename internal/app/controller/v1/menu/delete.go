package menu

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
)

func (ctrl *MenuController) delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := ctrl.biz.Menu().Delete(c, id); err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "删除成功", nil)
}

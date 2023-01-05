package role

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/log"
)

func (ctrl *RoleController) delete(c *gin.Context) {
	log.C(c.Request.Context()).Infow("delete Role function called")

	if err := ctrl.biz.Role().Delete(c.Request.Context(), c.Param("id")); err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "删除成功", nil)
}

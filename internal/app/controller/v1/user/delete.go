package user

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/log"
)

func (ctrl *UserController) delete(c *gin.Context) {
	log.C(c).Infow("Delete user function called")

	err := ctrl.biz.User().Delete(c, c.Param("username"))

	if err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "删除成功", nil)
}

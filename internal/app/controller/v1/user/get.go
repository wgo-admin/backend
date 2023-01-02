package user

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/log"
)

func (ctrl *UserController) get(c *gin.Context) {
	log.C(c).Infow("Get user function called")

	resp, err := ctrl.biz.User().Get(c, c.Param("username"))

	if err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "查询成功", resp)
}

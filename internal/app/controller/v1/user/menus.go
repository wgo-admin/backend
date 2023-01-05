package user

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/known"
	"github.com/wgo-admin/backend/internal/pkg/log"
)

func (ctrl *UserController) menus(c *gin.Context) {
	log.C(c.Request.Context()).Infow("User menus function called.")

	resp, err := ctrl.biz.Menu().UserMenus(c.Request.Context(), c.GetString(known.XRoleKey))
	if err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "查询成功", resp)
}

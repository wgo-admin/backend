package user

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/known"
	"github.com/wgo-admin/backend/internal/pkg/log"
)

func (ctrl *UserController) profile(c *gin.Context) {
	log.C(c.Request.Context()).Infow("User profile function called.")

	username := c.GetString(known.XUsernameKey)
	userInfo, err := ctrl.biz.User().Profile(c.Request.Context(), username)
	if err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "查询成功", userInfo)
}

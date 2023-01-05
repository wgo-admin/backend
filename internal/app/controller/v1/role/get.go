package role

import (
  "github.com/gin-gonic/gin"
  "github.com/wgo-admin/backend/internal/pkg/core"
  "github.com/wgo-admin/backend/internal/pkg/log"
)

func (ctrl *RoleController) get(c *gin.Context) {
  log.C(c.Request.Context()).Infow("Role get function called.")

  sysApi, err := ctrl.biz.Role().Get(c.Request.Context(), c.Param("id"))
  if err != nil {
    core.ResponseFail(c, err)
    return
  }

  core.ResponseOk(c, "查询成功", sysApi)
}

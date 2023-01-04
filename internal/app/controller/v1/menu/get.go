package menu

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
)

func (ctrl *MenuController) get(c *gin.Context) {
  id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

  resp, err := ctrl.biz.Menu().Get(c, id)
  if err != nil {
    core.ResponseFail(c, err)
    return
  }

  core.ResponseOk(c, "查询成功", resp)
}

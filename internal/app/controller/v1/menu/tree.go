package menu

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/log"
)

func (ctrl *MenuController) tree(c *gin.Context) {
	log.C(c).Infow("Menu tree function called")

	parentId, _ := strconv.ParseInt(c.Param("parentId"), 10, 64)

	resp, err := ctrl.biz.Menu().GetListByParentID(c, parentId)
	if err != nil {
		core.ResponseFail(c, err)
		return
	}

	core.ResponseOk(c, "查询成功", resp)
}

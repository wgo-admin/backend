package sysApi

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/app/biz"
	v1 "github.com/wgo-admin/backend/internal/app/controller/v1"
	mw "github.com/wgo-admin/backend/internal/pkg/middleware"
)

var name = "v1_sysApi"

var _ v1.IController = (*SysApiController)(nil)

func newController() *SysApiController {
	return &SysApiController{}
}

type SysApiController struct {
	biz biz.IBiz
}

func (ctrl *SysApiController) InjectBiz(biz biz.IBiz) {
	ctrl.biz = biz
}

func (ctrl *SysApiController) Name() string {
	return name
}

func (ctrl *SysApiController) RegistryApi(g gin.IRouter) {
	group := g.Group("/sys_apis", mw.Authn())
	{
		group.POST("", ctrl.create)
		group.GET("", ctrl.list)
		group.GET(":id", ctrl.get)
		group.PUT(":id", ctrl.update)
		group.DELETE("", ctrl.delete)
	}
}

func init() {
	ctrl := newController()
	v1.RegistryController(ctrl)
}

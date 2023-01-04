package menu

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/app/biz"
	v1 "github.com/wgo-admin/backend/internal/app/controller/v1"
	mw "github.com/wgo-admin/backend/internal/pkg/middleware"
)

var name = "v1_menu"

var _ v1.IController = (*MenuController)(nil)

func newController() *MenuController {
	return &MenuController{}
}

type MenuController struct {
	biz biz.IBiz
}

func (ctrl *MenuController) InjectBiz(biz biz.IBiz) {
	ctrl.biz = biz
}

func (ctrl *MenuController) Name() string {
	return name
}

func (ctrl *MenuController) RegistryApi(g gin.IRouter) {
	group := g.Group("/menus", mw.Authn())
	{
		group.POST("", ctrl.create)
		group.GET("/tree/:parentId", ctrl.tree)
	}
}

func init() {
	ctrl := newController()
	v1.RegistryController(ctrl)
}

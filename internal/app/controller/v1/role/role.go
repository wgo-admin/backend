package role

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/app/biz"
	v1 "github.com/wgo-admin/backend/internal/app/controller/v1"
	mw "github.com/wgo-admin/backend/internal/pkg/middleware"
)

var name = "v1_role"

var _ v1.IController = (*RoleController)(nil)

func newController() *RoleController {
	return &RoleController{}
}

type RoleController struct {
	biz biz.IBiz
}

func (ctrl *RoleController) InjectBiz(biz biz.IBiz) {
	ctrl.biz = biz
}

func (ctrl *RoleController) Name() string {
	return name
}

func (ctrl *RoleController) RegistryApi(g gin.IRouter) {
	group := g.Group("/roles", mw.Authn(), mw.Authz())
	{
		group.GET("", ctrl.list)
		group.GET(":id", ctrl.get)
		group.POST("", ctrl.create)
		group.PUT(":id", ctrl.update)
		group.DELETE(":id", ctrl.delete)
	}
}

func init() {
	ctrl := newController()
	v1.RegistryController(ctrl)
}

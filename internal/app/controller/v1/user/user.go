package user

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/app/biz"
	v1 "github.com/wgo-admin/backend/internal/app/controller/v1"
)

var name = "v1_user"

func newController() *UserController {
	return &UserController{}
}

type UserController struct {
	biz biz.IBiz
}

func (ctrl *UserController) InjectBiz(biz biz.IBiz) {
	ctrl.biz = biz
}

func (ctrl *UserController) Name() string {
	return name
}

func (ctrl *UserController) RegistryApi(g gin.IRouter) {
	group := g.Group("/users")
	{
		group.POST("/register", ctrl.register)
		group.POST("/login", ctrl.login)
	}
}

func init() {
	ctrl := newController()
	v1.RegistryController(ctrl)
}

package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/app/biz"
)

type IController interface {
	// RegistryApi 注册Api
	RegistryApi(g gin.IRouter)
	// InjectBiz 注入Biz层
	InjectBiz(biz biz.IBiz)
	// Name 名称
	Name() string
}

var (
	ctrls = map[string]IController{}
)

// RegistryController RegistryGin 每个模块可用该方法去注册api
func RegistryController(ctrl IController) {
	if _, ok := ctrls[ctrl.Name()]; ok {
		panic(any(fmt.Sprintf("controller %s has registried", ctrl.Name())))
	}
	ctrls[ctrl.Name()] = ctrl
}

// GetRegisteredControllers 获取已注册的 controller 的名称切片
func GetRegisteredControllers() (names []string) {
	for _, ctrl := range ctrls {
		names = append(names, ctrl.Name())
	}
	return
}

// InitControllers 初始化 controller
func InitControllers(g gin.IRouter, biz biz.IBiz) {
	// 创建 v1 分组
	v1 := g.Group("/v1")

	// 为每一个模块注入 biz 业务层
	for _, ctrl := range ctrls {
		ctrl.InjectBiz(biz)
	}

	// 注册每个模块的api
	for _, ctrl := range ctrls {
		ctrl.RegistryApi(v1)
	}
}

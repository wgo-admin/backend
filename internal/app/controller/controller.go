package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/app/biz"
	v1 "github.com/wgo-admin/backend/internal/app/controller/v1"
	"github.com/wgo-admin/backend/internal/pkg/log"

	_ "github.com/wgo-admin/backend/internal/app/controller/registry"
)

// Init 初始化 controller 层
func Init(g gin.IRouter, biz biz.IBiz) {
	// 注册 v1 版本的接口
	v1.InitControllers(g, biz)
	log.Infow("registry api v1", "controllers", v1.GetRegisteredControllers())
}

package app

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/app/biz"
	"github.com/wgo-admin/backend/internal/app/controller"
	"github.com/wgo-admin/backend/internal/app/store"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/log"
)

func registryRoutes(g *gin.Engine) error {
	// 注册 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		core.ResponseFail(c, errno.ErrPageNotFound)
	})

	// 注册 /health handler.
	g.GET("/health", func(c *gin.Context) {
		log.C(c.Request.Context()).Infow("Healthz function called")
		core.ResponseOk(c, "pong", map[string]string{"status": "ok"})
	})

	// 注册 pprof 路由
	pprof.Register(g)

	// 初始化 Biz 层 （业务层）
	biz := biz.NewBiz(store.S)

	// 初始化 controller 层
	controller.Init(g, biz)

	return nil
}

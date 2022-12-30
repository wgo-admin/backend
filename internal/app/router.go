package app

import (
  "github.com/gin-contrib/pprof"
  "github.com/gin-gonic/gin"
  "github.com/wgo-admin/backend/internal/pkg/errno"
  "github.com/wgo-admin/backend/internal/pkg/log"
  "github.com/wgo-admin/backend/internal/pkg/resp"
)

func registryRoutes(g *gin.Engine) error {
  // 注册 404 Handler.
  g.NoRoute(func(c *gin.Context) {
    resp.Write(c, errno.ErrPageNotFound, nil)
  })

  // 注册 /health handler.
  g.GET("/health", func(c *gin.Context) {
    log.C(c.Request.Context()).Infow("Healthz function called")

    resp.Write(c, nil, map[string]string{"status": "ok"})
  })

  // 注册 pprof 路由
  pprof.Register(g)

  return nil
}

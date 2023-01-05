package middleware

import (
  "github.com/gin-gonic/gin"
  "github.com/wgo-admin/backend/internal/app/store"
  "github.com/wgo-admin/backend/internal/pkg/core"
  "github.com/wgo-admin/backend/internal/pkg/errno"
  "github.com/wgo-admin/backend/internal/pkg/known"
)

// Authz 授权中间件
func Authz() gin.HandlerFunc {
  return func(c *gin.Context) {
    sub := getUserRoleId(c)
    obj := c.FullPath()
    act := c.Request.Method

    // 判断是否有权限
    success, _ := store.S.Casbins().Authorize(c.Request.Context(), sub, obj, act)
    if !success {
      core.WriteResponse(c, errno.ErrForbidden, nil)
      c.Abort()
      return
    }

    c.Next()
  }
}

func getUserRoleId(c *gin.Context) string {
  roleId := c.GetString(known.XRoleKey)
  return roleId
}

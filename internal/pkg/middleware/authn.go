package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wgo-admin/backend/internal/app/store"
	"github.com/wgo-admin/backend/internal/pkg/core"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/known"
	"github.com/wgo-admin/backend/pkg/token"
)

// Authn 认证用户中间件，如果合法则把 token 中的解析出来的 username 和 roleId存放到 gin.Context 中
func Authn() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, roleId, err := token.ParseRequest(c)
		if err != nil {
			core.WriteResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		// 从 redis 中获取用户是否登录
		_, err = store.S.Cache().GetStruct(c.Request.Context(), known.GetRedisUserKey(username))
		if err != nil {
			core.WriteResponse(c, errno.ErrUserLogout, nil)
			c.Abort()
			return
		}

		c.Set(known.XUsernameKey, username)
		c.Set(known.XRoleKey, roleId)
		c.Next()
	}
}

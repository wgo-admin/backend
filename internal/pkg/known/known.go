package known

import "fmt"

const (
  // XRequestIDKey 用来定义 Gin 上下文中的键，代表请求的 uuid.
  XRequestIDKey = "X-Request-ID"

  // XUsernameKey 用来定义 Gin 上下文的键，代表请求的所有者.
  XUsernameKey = "X-Username"

  // XRoleKey 用来定义 Gin 上下文的键，代表请求的所有者的角色ID.
  XRoleKey = "X-Role-ID"
)

func GetRedisUserKey(username string) string {
  return fmt.Sprintf("wgo:user:%s", username)
}

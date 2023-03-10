// 注册每个模块的 controller
// 因为每个模块都需要有个 init 方法去注册 controller 到 IOC 里
// 所以需要引包，去调用每个模块的 init 方法

package registry

import (
  _ "github.com/wgo-admin/backend/internal/app/controller/v1/menu"
  _ "github.com/wgo-admin/backend/internal/app/controller/v1/role"
  _ "github.com/wgo-admin/backend/internal/app/controller/v1/sysApi"
  _ "github.com/wgo-admin/backend/internal/app/controller/v1/user"
)

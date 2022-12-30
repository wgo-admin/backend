package auth

import (
	"time"

	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const (
	// casbin 访问控制模型.
	rbacModel = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")`
)

// Authz 定义了一个授权器，提供授权功能.
type Authz struct {
	*casbin.SyncedEnforcer
}

// NewAuthz 创建一个使用 casbin 完成授权的授权器.
func NewAuthz(db *gorm.DB) (*Authz, error) {
	// Initialize a Gorm adapter and use it in a Casbin enforcer
	adp, err := adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 加载rbac权限模型
	m, _ := model.NewModelFromString(rbacModel)

	// Initialize the enforcer.
	enforcer, err := casbin.NewSyncedEnforcer(m, adp)
	if err != nil {
		return nil, err
	}

	// Load the policy from DB.
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	enforcer.StartAutoLoadPolicy(5 * time.Second)

	a := &Authz{enforcer}

	return a, nil
}

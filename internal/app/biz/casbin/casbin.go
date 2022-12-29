package casbin

import (
	"context"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/wgo-admin/backend/internal/app/store"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
	"github.com/wgo-admin/backend/pkg/auth"
)

type ICasbinBiz interface {
	Clear(ctx context.Context, fieldIndex int, p ...string) bool
	CreateOrUpdate(ctx context.Context, roleId string, infos []v1.PermissionInfo) error
	UpdateRule(ctx context.Context, oldPath, oldMethod, newPath, newMethod string) error
}

type casbinBiz struct {
	ds    store.IStore
	authz *auth.Authz
}

func New(ds store.IStore, authz *auth.Authz) *casbinBiz {
	return &casbinBiz{ds, authz}
}

var _ ICasbinBiz = (*casbinBiz)(nil)

// 更新了某一个接口则也需要更新casbin里的rule
func (c *casbinBiz) UpdateRule(ctx context.Context, oldPath, oldMethod, newPath, newMethod string) error {
	return c.ds.DB().
		Model(&gormadapter.CasbinRule{}).
		Where("v1 = ? AND v2 = ?", oldPath, oldMethod).
		Updates(map[string]string{
			"v1": newPath,
			"v2": newMethod,
		}).Error
}

// 给角色创建或更新权限
func (c *casbinBiz) CreateOrUpdate(ctx context.Context, roleId string, infos []v1.PermissionInfo) error {
	c.Clear(ctx, 0, roleId)
	rules := [][]string{}
	for _, v := range infos {
		rules = append(rules, []string{roleId, v.Path, v.Method})
	}
	success, _ := c.authz.AddPolicy(rules)
	if !success {
		return errno.ErrHasSameApi
	}
	return nil
}

// 清除匹配的权限
func (c *casbinBiz) Clear(ctx context.Context, fieldIndex int, p ...string) bool {
	success, _ := c.authz.RemoveFilteredPolicy(fieldIndex, p...)
	return success
}

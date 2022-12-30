package store

import (
	"context"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/pkg/auth"
	"gorm.io/gorm"
)

type ICasbinStore interface {
	Clear(ctx context.Context, fieldIndex int, p ...string) bool
	CreateOrUpdate(ctx context.Context, roleId string, infos []map[string]string) error
	UpdateRule(ctx context.Context, oldPath, oldMethod, newPath, newMethod string) error
	Authorize(ctx context.Context, sub, obj, act string) (bool, error)
}

var _ ICasbinStore = (*casbins)(nil)

func newCasbins(db *gorm.DB, authz *auth.Authz) *casbins {
	return &casbins{
		db:     db,
		casbin: authz,
	}
}

type casbins struct {
	db     *gorm.DB
	casbin *auth.Authz
}

// 清除匹配的权限
func (s *casbins) Clear(ctx context.Context, fieldIndex int, p ...string) bool {
	success, _ := s.casbin.RemoveFilteredPolicy(fieldIndex, p...)
	return success
}

// 给角色创建或更新权限
func (s *casbins) CreateOrUpdate(ctx context.Context, roleId string, infos []map[string]string) error {
	s.Clear(ctx, 0, roleId)
	var rules [][]string
	for _, v := range infos {
		rules = append(rules, []string{roleId, v["Path"], v["Method"]})
	}
	success, _ := s.casbin.AddPolicy(rules)
	if !success {
		return errno.ErrHasSameApi
	}
	return nil
}

// 更新了某一个接口则也需要更新casbin里的rule
func (s *casbins) UpdateRule(ctx context.Context, oldPath, oldMethod, newPath, newMethod string) error {
	return s.db.
		Model(&gormadapter.CasbinRule{}).
		Where("v1 = ? AND v2 = ?", oldPath, oldMethod).
		Updates(map[string]string{
			"v1": newPath,
			"v2": newMethod,
		}).Error
}

// Authorize 用来进行授权.
func (s *casbins) Authorize(ctx context.Context, sub, obj, act string) (bool, error) {
	return s.casbin.Enforce(sub, obj, act)
}

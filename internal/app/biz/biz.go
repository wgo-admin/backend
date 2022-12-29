package biz

import (
	"github.com/wgo-admin/backend/internal/app/biz/casbin"
	"github.com/wgo-admin/backend/internal/app/store"
	"github.com/wgo-admin/backend/pkg/auth"
)

// 业务层接口
type IBiz interface {
	Casbin() casbin.ICasbinBiz
}

var _ IBiz = (*biz)(nil)

// 创建 biz 实例
func NewBiz(ds store.IStore, authz *auth.Authz) *biz {
	return &biz{
		ds:    ds,
		authz: authz,
	}
}

type biz struct {
	ds    store.IStore
	authz *auth.Authz
}

func (b *biz) Casbin() casbin.ICasbinBiz {
	return casbin.New(b.ds, b.authz)
}

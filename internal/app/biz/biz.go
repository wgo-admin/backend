package biz

import (
	"github.com/wgo-admin/backend/internal/app/biz/sysApi"
	"github.com/wgo-admin/backend/internal/app/biz/user"
	"github.com/wgo-admin/backend/internal/app/store"
)

// 业务层接口
type IBiz interface {
	SysApi() sysApi.ISysApiBiz
	User() user.IUserBiz
}

var _ IBiz = (*biz)(nil)

// 创建 biz 实例
func NewBiz(ds store.IStore) *biz {
	return &biz{
		ds: ds,
	}
}

type biz struct {
	ds store.IStore
}

func (b *biz) User() user.IUserBiz {
	return user.NewBiz(b.ds)
}

func (b *biz) SysApi() sysApi.ISysApiBiz {
	return sysApi.NewBiz(b.ds)
}

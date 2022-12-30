package store

import (
	"github.com/go-redis/redis/v8"
	"github.com/wgo-admin/backend/pkg/auth"
	"sync"

	"github.com/wgo-admin/backend/internal/app/store/model"
	"gorm.io/gorm"
)

type IStore interface {
	DB() *gorm.DB
	Casbins() ICasbinStore
}

var (
	once sync.Once
	// 全局变量，方便其它包直接调用已初始化好的 S 实例.
	S *datastore
)

var _ IStore = (*datastore)(nil)

func NewStore(db *gorm.DB, rdb *redis.Client, casbin *auth.Authz) *datastore {
	// 确保只执行一次
	once.Do(func() {
		S = &datastore{db: db, rdb: rdb, casbin: casbin}
	})
	return S
}

type datastore struct {
	db     *gorm.DB
	rdb    *redis.Client
	casbin *auth.Authz
}

// Permissions 返回 IPermissionsStore 接口实例
func (ds *datastore) Casbins() ICasbinStore {
	return NewCasbins(ds.db, ds.casbin)
}

// 返回一个 gorm db 对象
func (ds *datastore) DB() *gorm.DB {
	return ds.db
}

// 将模型迁移到数据库
func (ds *datastore) AutoMigrate() error {
	return ds.db.AutoMigrate(model.RoleM{}, model.UserM{}, model.MenuM{}, model.SysApiM{})
}

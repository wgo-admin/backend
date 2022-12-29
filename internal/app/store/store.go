package store

import (
	"github.com/go-redis/redis/v8"
	"sync"

	"github.com/wgo-admin/backend/internal/app/store/model"
	"gorm.io/gorm"
)

var (
	once sync.Once
	// 全局变量，方便其它包直接调用已初始化好的 S 实例.
	S *datastore
)

func NewStore(db *gorm.DB, rdb *redis.Client) *datastore {
	// 确保只执行一次
	once.Do(func() {
		S = &datastore{db: db, rdb: rdb}
	})
	return S
}

type datastore struct {
	db  *gorm.DB
	rdb *redis.Client
}

// 返回一个 gorm db 对象
func (ds *datastore) DB() *gorm.DB {
	return ds.db
}

func (ds *datastore) AutoMigrate() error {
	return ds.db.AutoMigrate(model.RoleM{}, model.UserM{}, model.MenuM{}, model.SysApiM{})
}

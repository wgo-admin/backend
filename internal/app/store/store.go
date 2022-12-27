package store

import (
	"sync"

	"gorm.io/gorm"
)

var (
	once sync.Once
	// 全局变量，方便其它包直接调用已初始化好的 S 实例.
	S *datastore
)

func NewStore(db *gorm.DB) *datastore {
	// 确保只执行一次
	once.Do(func() {
		S = &datastore{db: db}
	})
	return S
}

type datastore struct {
	db *gorm.DB
}

// 返回一个 gorm db 对象
func (ds *datastore) DB() *gorm.DB {
	return ds.db
}

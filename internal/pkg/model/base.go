package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseM 基础模型
type BaseM struct {
	ID        int64          `gorm:"column:id;primaryKey"`    // 主键ID
	CreatedAt time.Time      `gorm:"column:created_at"`       // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at"`       // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"` // 删除时间
}

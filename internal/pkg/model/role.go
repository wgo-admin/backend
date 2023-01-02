package model

import (
	"time"

	"gorm.io/gorm"
)

type RoleM struct {
	ID        string         `gorm:"column:id;not null;primaryKey;comment:角色id;size:90"` // 角色id
	CreatedAt time.Time      `gorm:"column:created_at"`                                  // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at"`                                  // 创建时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`                            // 删除时间
	Name      string         `gorm:"column:name;not null;index;comment:角色名称"`
	Desc      string         `gorm:"column:desc;comment:角色描述"`
	UserMs    []UserM        `gorm:"foreignKey:RoleID;references:ID"`                                                                 // 关联外键是 UserM 的 RoleID 字段，跟 RoleID 字段相关联
	MenuMs    *[]MenuM       `gorm:"many2many:role_menus;foreignKey:ID;joinForeignKey:role_id;references:ID;joinReferences:menu_id;"` // 角色菜单多对多关系
}

func (r *RoleM) TableName() string {
	return "role"
}

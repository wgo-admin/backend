package model

import (
	"github.com/wgo-admin/backend/pkg/auth"
	"gorm.io/gorm"
)

type UserM struct {
	BaseM
	Username string `gorm:"column:username;unique;not null;comment:用户名"`
	Password string `gorm:"column:password;not null;comment:登录密码"`
	Nickname string `gorm:"column:nickname;comment:昵称"`
	Email    string `gorm:"column:email;unique;not null;comment:邮箱"`
	Phone    string `gorm:"column:phone;comment:电话号码"`
	RoleID   string
}

// TableName 用来指定映射的 MySQL 表名.
func (u *UserM) TableName() string {
	return "user"
}

// BeforeCreate 在创建数据库记录之前加密明文密码.
func (u *UserM) BeforeCreate(tx *gorm.DB) (err error) {
	// Encrypt the user password.
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
	}

	return nil
}

package store

import (
	"context"
	"errors"

	"github.com/wgo-admin/backend/internal/app/store/model"
	"gorm.io/gorm"
)

type IUserStore interface {
	Create(ctx context.Context, user *model.UserM) error
	Get(ctx context.Context, username string) (*model.UserM, error)
	Update(ctx context.Context, user *model.UserM) error
	Delete(ctx context.Context, username string) error
	List(ctx context.Context, user *model.UserM, offset, limit int) (int64, []*model.UserM, error)
}

var _ IUserStore = (*users)(nil)

func newUsers(db *gorm.DB) *users {
	return &users{db}
}

type users struct {
	db *gorm.DB
}

// Create 创建用户
func (u *users) Create(ctx context.Context, user *model.UserM) error {
	return u.db.Create(&user).Error
}

// Get 获取用户详情
func (u *users) Get(ctx context.Context, username string) (*model.UserM, error) {
	var user model.UserM
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (u *users) Update(ctx context.Context, user *model.UserM) error {
	return u.db.Save(&user).Error
}

// Delete 删除用户
func (u *users) Delete(ctx context.Context, username string) error {
	err := u.db.Where("username = ?", username).Delete(&model.UserM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

// List 查询列表
func (u *users) List(ctx context.Context, user *model.UserM, offset, limit int) (total int64, ret []*model.UserM, err error) {
	// 添加条件查询
	db := u.db.Model(&model.UserM{})
	if user.Username != "" {
		db.Where("username LIKE ?", "%"+user.Username+"%")
	}
	if user.Email != "" {
		db.Where("email LIKE ?", "%"+user.Email+"%")
	}

	// 条件过滤后的数据总数
	err = db.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	// 分页查询
	if err := db.Offset(offset).Limit(limit).Order("id desc").Find(&ret).Error; err != nil {
		return 0, nil, err
	}

	return
}

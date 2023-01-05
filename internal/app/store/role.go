package store

import (
	"context"
	"errors"
	"github.com/wgo-admin/backend/internal/pkg/log"

	"github.com/wgo-admin/backend/internal/pkg/model"
	"gorm.io/gorm"
)

type IRoleStore interface {
	Create(ctx context.Context, role *model.RoleM) error
	Get(ctx context.Context, id string) (*model.RoleM, error)
	Update(ctx context.Context, role *model.RoleM) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, role *model.RoleM, offset, limit int) (total int64, ret []*model.RoleM, err error)
}

var _ IRoleStore = (*roles)(nil)

func newRoles(db *gorm.DB) *roles {
	return &roles{db}
}

type roles struct {
	db *gorm.DB
}

func (r *roles) Create(ctx context.Context, role *model.RoleM) error {
	return r.db.Create(&role).Error
}

func (r *roles) Delete(ctx context.Context, id string) error {
	var roleM model.RoleM
	err := r.db.Where("id = ?", id).First(&roleM).Delete(&roleM).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 解除与菜单的关系
	if err := r.db.Model(&roleM).Association("MenusM").Clear(); err != nil {
		log.C(ctx).Errorw("RoleM clear MenusM relation failed", "error", err)
		return err
	}

	return nil
}

func (r *roles) Get(ctx context.Context, id string) (*model.RoleM, error) {
	var role model.RoleM
	if err := r.db.Preload("UsersM").Preload("MenusM").First(&role, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roles) List(ctx context.Context, role *model.RoleM, offset int, limit int) (total int64, ret []*model.RoleM, err error) {
	// 添加条件查询
	db := r.db.Preload("UsersM").Preload("MenusM").Model(&model.RoleM{})
	if role.ID != "" {
		db.Where("id LIKE ?", "%"+role.ID+"%")
	}
	if role.Name != "" {
		db.Where("name LIKE ?", "%"+role.Name+"%")
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

func (r *roles) Update(ctx context.Context, role *model.RoleM) error {
	return r.db.Save(&role).Error
}

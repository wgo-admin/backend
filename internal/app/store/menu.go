package store

import (
	"context"
	"errors"

	"github.com/wgo-admin/backend/internal/pkg/log"
	"github.com/wgo-admin/backend/internal/pkg/model"
	"gorm.io/gorm"
)

type IMenuStore interface {
	Create(ctx context.Context, menu *model.MenuM) error
	Get(ctx context.Context, id int64) (*model.MenuM, error)
	Update(ctx context.Context, menu *model.MenuM) error
	Delete(ctx context.Context, ids int64) error
	All(ctx context.Context) (ret []*model.MenuM, err error)
}

var _ IMenuStore = (*menus)(nil)

func newMenus(db *gorm.DB) *menus {
	return &menus{db}
}

type menus struct {
	db *gorm.DB
}

// 创建菜单
func (m *menus) Create(ctx context.Context, menu *model.MenuM) error {
	return m.db.Create(&menu).Error
}

// 获取菜单详情
func (m *menus) Get(ctx context.Context, id int64) (*model.MenuM, error) {
	var menu model.MenuM
	if err := m.db.Preload("SysApisM").First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// 更新菜单
func (m *menus) Update(ctx context.Context, menu *model.MenuM) error {
	return m.db.Save(&menu).Error
}

// 删除菜单
func (m *menus) Delete(ctx context.Context, ids int64) error {
  var menuM model.MenuM
	err := m.db.Preload("SysApisM").Where("id = ?", ids).Find(&menuM).Delete(&model.MenuM{}).Error
	// 如果有错误且错误不是 ErrRecordNotFound，则返回err
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

  // 解除与接口的关系
  if err := m.db.Model(&menuM).Association("SysApisM").Delete(&menuM.SysApisM); err != nil {
    log.C(ctx).Errorw("MenuM unbind SysApi relation failed", "error", err)
    return err
  }

	return nil
}

// 获取所有的菜单
func (m *menus) All(ctx context.Context) (ret []*model.MenuM, err error) {
	err = m.db.Model(&model.MenuM{}).Find(&ret).Error
	return
}

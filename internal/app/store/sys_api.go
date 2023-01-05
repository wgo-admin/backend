package store

import (
	"context"
	"github.com/wgo-admin/backend/pkg/auth"

	"github.com/wgo-admin/backend/internal/pkg/log"
	"github.com/wgo-admin/backend/internal/pkg/model"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type ISysApiStore interface {
	Create(ctx context.Context, sysApi *model.SysApiM) error
	Get(ctx context.Context, id int64) (*model.SysApiM, error)
	Update(ctx context.Context, sysApi *model.SysApiM) error
	Delete(ctx context.Context, ids []int64) error
	List(ctx context.Context, sysApi *model.SysApiM, offset, limit int) (int64, []*model.SysApiM, error)
}

type sysApis struct {
	db     *gorm.DB
	casbin *auth.Authz
}

var _ ISysApiStore = (*sysApis)(nil)

func newSysApis(db *gorm.DB, casbin *auth.Authz) *sysApis {
	return &sysApis{db: db, casbin: casbin}
}

func (s *sysApis) Create(ctx context.Context, sysApi *model.SysApiM) error {
	return s.db.Create(&sysApi).Error
}

func (s *sysApis) Get(ctx context.Context, id int64) (*model.SysApiM, error) {
	var sysApi model.SysApiM
	if err := s.db.First(&sysApi, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &sysApi, nil
}

func (s *sysApis) Update(ctx context.Context, sysApi *model.SysApiM) error {
	return s.db.Save(sysApi).Error
}

func (s *sysApis) Delete(ctx context.Context, ids []int64) error {
	var sysApisM []*model.SysApiM

	if err := s.db.Preload("MenusM").Where("id in (?)", ids).Find(&sysApisM).Delete(&sysApisM).Error; err != nil {
		return err
	}

	// 使用 goroutine 提升性能
	eg, ctx := errgroup.WithContext(ctx)
	for _, item := range sysApisM {
		sysApiM := item
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				if len(sysApiM.MenusM) > 0 {
					// 解除与菜单之间的关系
					if err := s.db.Model(&sysApiM).Association("MenusM").Delete(&sysApiM.MenusM); err != nil {
						log.C(ctx).Errorw("SysApi unbind menu relation failed", "error", err)
						return err
					}
				}

				// 同时删除 casbin_rule 表的策略规则
				_, err := s.casbin.RemoveFilteredPolicy(1, sysApiM.Path, sysApiM.Method)
				if err != nil {
					return err
				}

				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.C(ctx).Errorw("Failed to wait all function calls returned", "err", err)
		return err
	}

	return nil
}

func (s *sysApis) List(ctx context.Context, sysApi *model.SysApiM, offset, limit int) (total int64, ret []*model.SysApiM, err error) {
	db := s.db.Model(&model.SysApiM{})

	if sysApi.Method != "" {
		db = db.Where("method LIKE ?", "%"+sysApi.Method+"%")
	}

	if sysApi.GroupName != "" {
		db = db.Where("group_name LIKE ?", "%"+sysApi.GroupName+"%")
	}

	if sysApi.Path != "" {
		db = db.Where("path LIKE ?", "%"+sysApi.Path+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	db = db.Offset(offset).Limit(limit).Order("id desc")
	err = db.Find(&ret).Error
	if err != nil {
		return 0, nil, err
	}

	return
}

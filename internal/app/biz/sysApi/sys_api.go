package sysApi

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
	"github.com/wgo-admin/backend/internal/app/store"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/log"
	"github.com/wgo-admin/backend/internal/pkg/model"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
	"gorm.io/gorm"
)

type ISysApiBiz interface {
	Create(ctx context.Context, req *v1.CreateSysApiRequest) error
	List(ctx context.Context, req *v1.QuerySysApiRequest) (resp *v1.QuerySysApiResponse, err error)
	BatchDelete(ctx context.Context, req *v1.IdsRequest) error
	Update(ctx context.Context, id int64, req *v1.UpdateSysApiRequest) error
	Get(ctx context.Context, id int64) (*v1.SysApiInfo, error)
}

var _ ISysApiBiz = (*sysApiBiz)(nil)

func NewBiz(ds store.IStore) *sysApiBiz {
	return &sysApiBiz{ds}
}

type sysApiBiz struct {
	ds store.IStore
}

// 批量删除
func (b *sysApiBiz) BatchDelete(ctx context.Context, req *v1.IdsRequest) error {
	return b.ds.SysApis().Delete(ctx, *req)
}

// 获取详细信息
func (b *sysApiBiz) Get(ctx context.Context, id int64) (*v1.SysApiInfo, error) {
	sysApiM, err := b.ds.SysApis().Get(ctx, id)
	if err != nil {
		return nil, err
	}
	sysApiInfo := v1.SysApiInfo{}
	copier.Copy(&sysApiInfo, sysApiM)
	sysApiInfo.CreatedAt = sysApiM.CreatedAt.Format(v1.TIME_FORMAT)
	sysApiInfo.UpdatedAt = sysApiM.UpdatedAt.Format(v1.TIME_FORMAT)

	return &sysApiInfo, nil
}

// 更新
func (b *sysApiBiz) Update(ctx context.Context, id int64, req *v1.UpdateSysApiRequest) error {
	sysApiM, err := b.ds.SysApis().Get(ctx, id)
	if err != nil {
		return err
	}

  oldPath := sysApiM.Path
  oldMethod := sysApiM.Method

	// 需要查询有没有相同的method和path，有则不能去更新
	err = b.ds.DB().Where("method = ? AND path = ?", req.Method, req.Path).First(&sysApiM).Error
	// 如果错误不是 ErrRecordNotFound，则需要返回错误, 找不到才可以去更新
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errno.ErrSysApiAlreadyExist
	}

	if req.Title != nil {
		sysApiM.Title = *req.Title
	}

	if req.Method != nil {
		sysApiM.Method = *req.Method
	}

	if req.Path != nil {
		sysApiM.Path = *req.Path
	}

	if req.GroupName != nil {
		sysApiM.GroupName = *req.GroupName
	}

	err = b.ds.SysApis().Update(ctx, sysApiM)
	if err != nil {
		return err
	}

  // 更新完api还需更新casbin的策略规则
  if err := b.ds.Casbins().UpdateRule(ctx, oldPath, oldMethod, *req.Path, *req.Method); err != nil {
    log.C(ctx).Errorw("casbin UpdateRule failed", "error", err)
    return err
  }

	return nil
}

// 查询列表
func (b *sysApiBiz) List(ctx context.Context, req *v1.QuerySysApiRequest) (resp *v1.QuerySysApiResponse, err error) {
	var sysApiM model.SysApiM
	copier.Copy(&sysApiM, req)
	total, list, err := b.ds.SysApis().List(ctx, &sysApiM, req.Offset(), req.Limit())
	if err != nil {
		log.C(ctx).Errorw("Faild to list sys_apis from storage", "err", err)
		return nil, err
	}

	sysApis := make([]*v1.SysApiInfo, 0, len(list))
	for _, item := range list {
		sysApis = append(sysApis, &v1.SysApiInfo{
			ID:        item.ID,
			CreatedAt: item.CreatedAt.Format(v1.TIME_FORMAT),
			UpdatedAt: item.UpdatedAt.Format(v1.TIME_FORMAT),
			Method:    item.Method,
			Path:      item.Path,
			GroupName: item.GroupName,
			Title:     item.Title,
		})
	}
	return &v1.QuerySysApiResponse{Total: total, List: sysApis}, nil
}

// 创建
func (b *sysApiBiz) Create(ctx context.Context, req *v1.CreateSysApiRequest) error {
	var sysApiM model.SysApiM

	if err := b.ds.DB().Where("method = ? AND path = ?", req.Method, req.Path).First(&sysApiM).Error; err != nil {
		// 如果错误不是ErrRecordNotFound，则返回错误
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ErrSysApiAlreadyExist
		}
	}

	copier.Copy(&sysApiM, req)

	if err := b.ds.SysApis().Create(ctx, &sysApiM); err != nil {
		return err
	}

	return nil
}

package menu

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/wgo-admin/backend/internal/app/store"
	"github.com/wgo-admin/backend/internal/pkg/log"
	"github.com/wgo-admin/backend/internal/pkg/model"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
	"golang.org/x/sync/errgroup"
)

type IMenuBiz interface {
	Create(ctx context.Context, req *v1.CreateOrUpdateMenuRequest) error
	GetListByParentID(ctx context.Context, parentId int64) (*v1.QueryMenuTreeResponse, error)
}

var _ IMenuBiz = (*menuBiz)(nil)

func NewBiz(ds store.IStore) *menuBiz {
	return &menuBiz{ds}
}

type menuBiz struct {
	ds store.IStore
}

func (b *menuBiz) GetListByParentID(ctx context.Context, parentId int64) (*v1.QueryMenuTreeResponse, error) {
	var menusM []*model.MenuM
	if err := b.ds.DB().Preload("SysApisM").Where("parent_id = ?", parentId).Find(&menusM).Error; err != nil {
		return nil, err
	}

	list := make([]*v1.MenuInfo, 0, len(menusM))

	// 获取绑定的 apiIds
	eg, ctx := errgroup.WithContext(ctx)

	// 序列化数据
	for _, item := range menusM {
		menu := item
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				var menus []*model.MenuM
				if err := b.ds.DB().Where("parent_id = ?", menu.ID).Find(&menus).Error; err != nil {
					return err
				}

				// 获取关联的api的id切片
				apiIds := []int64{}
				if len(menu.SysApisM) > 0 {
					for _, api := range menu.SysApisM {
						apiIds = append(apiIds, api.ID)
					}
				}

				// 是否是叶子结点
				isLeaf := len(menus) == 0

				list = append(list, &v1.MenuInfo{
					ID:        menu.ID,
					CreatedAt: menu.CreatedAt.Format(v1.TIME_FORMAT),
					UpdatedAt: menu.UpdatedAt.Format(v1.TIME_FORMAT),
					ParentID:  menu.ParentID,
					Title:     menu.Title,
					Type:      menu.Type,
					Sort:      menu.Sort,
					Icon:      menu.Icon,
					Component: menu.Component,
					Path:      menu.Path,
					Perm:      menu.Perm,
					Hidden:    menu.Hidden,
					IsLink:    menu.IsLink,
					KeepAlive: menu.KeepAlive,
					IsLeaf:    isLeaf,
					ApiIds:    apiIds,
				})
				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.C(ctx).Errorw("Failed to wait all function calls returned", "err", err)
		return nil, err
	}

	return &v1.QueryMenuTreeResponse{List: list}, nil
}

// 创建 menu
func (b *menuBiz) Create(ctx context.Context, req *v1.CreateOrUpdateMenuRequest) error {
	var menuM model.MenuM

	// 查找需要关联的 api
	var sysApisM []model.SysApiM
	if err := b.ds.DB().Where("id in (?)", req.ApiIds).Find(&sysApisM).Error; err != nil {
		return err
	}

	// 创建菜单
	copier.Copy(&menuM, req)
	menuM.SysApisM = sysApisM
	if err := b.ds.Menus().Create(ctx, &menuM); err != nil {
		return err
	}

	return nil
}

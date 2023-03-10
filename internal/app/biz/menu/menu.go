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
  Get(ctx context.Context, id int64) (*v1.MenuInfo, error)
  Delete(ctx context.Context, id int64) error
  Update(ctx context.Context, id int64, req *v1.CreateOrUpdateMenuRequest) error
  UserMenus(ctx context.Context, roleId string) (*v1.QueryUserMenuResponse, error)
}

var _ IMenuBiz = (*menuBiz)(nil)

func NewBiz(ds store.IStore) *menuBiz {
  return &menuBiz{ds}
}

type menuBiz struct {
  ds store.IStore
}

// 查询用户的菜单
func (b *menuBiz) UserMenus(ctx context.Context, roleId string) (*v1.QueryUserMenuResponse, error) {
  roleM, err := b.ds.Roles().Get(ctx, roleId)
  if err != nil {
    return nil, err
  }

  menus := make([]*v1.MenuInfo, 0, len(roleM.MenusM))
  for _, item := range roleM.MenusM {
    menu := item
    // 过滤掉按钮类型
    if menu.Type != "B" {
      menus = append(menus, &v1.MenuInfo{
        ID:        menu.ID,
        CreatedAt: menu.CreatedAt.Format(v1.TIME_FORMAT),
        UpdatedAt: menu.UpdatedAt.Format(v1.TIME_FORMAT),
        ParentID:  menu.ParentID,
        Title:     menu.Title,
        Sort:      menu.Sort,
        Type:      menu.Type,
        Icon:      menu.Icon,
        Component: menu.Component,
        Path:      menu.Path,
        Perm:      menu.Perm,
        Hidden:    menu.Hidden,
        IsLink:    menu.IsLink,
        KeepAlive: menu.KeepAlive,
      })
    }
  }

  return &v1.QueryUserMenuResponse{List: menus}, nil
}

// 更新
func (b *menuBiz) Update(ctx context.Context, id int64, req *v1.CreateOrUpdateMenuRequest) error {
  menuM, err := b.ds.Menus().Get(ctx, id)
  if err != nil {
    return err
  }
  copier.Copy(&menuM, req)

  // 查找要绑定的api
  var sysApisM []model.SysApiM
  if err := b.ds.DB().Where("id in (?)", req.ApiIds).Find(&sysApisM).Error; err != nil {
    return err
  }
  // 替换绑定关系
  menuM.SysApisM = sysApisM

  if err := b.ds.Menus().Update(ctx, menuM); err != nil {
    return err
  }

  return nil
}

// 删除菜单
func (b *menuBiz) Delete(ctx context.Context, id int64) error {
  if err := b.ds.Menus().Delete(ctx, id); err != nil {
    return err
  }

  return nil
}

// 获取详细信息
func (b *menuBiz) Get(ctx context.Context, id int64) (*v1.MenuInfo, error) {
  menuM, err := b.ds.Menus().Get(ctx, id)
  if err != nil {
    return nil, err
  }

  sysApiLen := len(menuM.SysApisM)
  sysApiIds := make([]int64, 0, sysApiLen)
  if sysApiLen > 0 {
    for _, item := range menuM.SysApisM {
      sysApi := item
      sysApiIds = append(sysApiIds, sysApi.ID)
    }
  }

  var menuInfo v1.MenuInfo
  copier.Copy(&menuInfo, menuM)
  menuInfo.CreatedAt = menuM.CreatedAt.Format(v1.TIME_FORMAT)
  menuInfo.UpdatedAt = menuM.UpdatedAt.Format(v1.TIME_FORMAT)
  menuInfo.ApiIds = sysApiIds

  return &menuInfo, nil
}

// 根据 parent_id 查询菜单列表
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

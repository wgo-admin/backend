package role

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/wgo-admin/backend/internal/app/store"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/model"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
	"gorm.io/gorm"
	"regexp"
)

type IRoleBiz interface {
	Create(ctx context.Context, req *v1.CreateRoleRequest) error
	Update(ctx context.Context, id string, req *v1.UpdateRoleRequest) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*v1.RoleInfo, error)
	List(ctx context.Context, req *v1.QueryRoleRequest) (*v1.QueryRoleResponse, error)
}

var _ IRoleBiz = (*roleBiz)(nil)

func NewBiz(ds store.IStore) *roleBiz {
	return &roleBiz{ds}
}

type roleBiz struct {
	ds store.IStore
}

func (b *roleBiz) List(ctx context.Context, req *v1.QueryRoleRequest) (*v1.QueryRoleResponse, error) {
	var roleM model.RoleM
	_ = copier.Copy(&roleM, req)
	total, list, err := b.ds.Roles().List(ctx, &roleM, req.Offset(), req.Limit())
	if err != nil {
		return nil, err
	}

	roles := make([]*v1.RoleInfo, 0, len(list))
	for _, item := range list {
		role := item

		menuIds := make([]int64, 0, len(role.MenusM))
		if len(role.MenusM) > 0 {
			for _, menu := range role.MenusM {
				menuIds = append(menuIds, menu.ID)
			}
		}

		roles = append(roles, &v1.RoleInfo{
			ID:        role.ID,
			CreatedAt: role.CreatedAt.Format(v1.TIME_FORMAT),
			UpdatedAt: role.UpdatedAt.Format(v1.TIME_FORMAT),
			Name:      role.Name,
			Desc:      role.Desc,
			UserCount: len(role.UsersM),
			MenuIds:   menuIds,
		})
	}

	return &v1.QueryRoleResponse{
		Total: total,
		List:  roles,
	}, nil
}

func (b *roleBiz) Get(ctx context.Context, id string) (*v1.RoleInfo, error) {
	roleM, err := b.ds.Roles().Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// 获取绑定的菜单id
	menuIds := make([]int64, 0, len(roleM.MenusM))
	if len(roleM.MenusM) > 0 {
		for _, menu := range roleM.MenusM {
			menuIds = append(menuIds, menu.ID)
		}
	}

	userCount := len(roleM.UsersM)
	var roleInfo v1.RoleInfo
	_ = copier.Copy(&roleInfo, roleM)
	roleInfo.MenuIds = menuIds
	roleInfo.UserCount = userCount
	roleInfo.CreatedAt = roleM.CreatedAt.Format(v1.TIME_FORMAT)
	roleInfo.UpdatedAt = roleM.UpdatedAt.Format(v1.TIME_FORMAT)

	return &roleInfo, nil
}

func (b *roleBiz) Delete(ctx context.Context, id string) error {
	err := b.ds.Roles().Delete(ctx, id)
	if err != nil {
		return err
	}

	// 删除角色时同时清除角色的权限
	if ok := b.ds.Casbins().Clear(ctx, 0, id); !ok {
		return errno.ErrClearPermissionFailed
	}

	return nil
}

func (b *roleBiz) Update(ctx context.Context, id string, req *v1.UpdateRoleRequest) error {
	// 获取角色详情
	roleM, err := b.ds.Roles().Get(ctx, id)
	if err != nil {
		return err
	}
	_ = copier.Copy(&roleM, req)

	// 替换与菜单的关系
	var menus []*model.MenuM
	if err := b.ds.DB().Preload("SysApisM").Where("id in (?)", req.MenuIds).Find(&menus).Error; err != nil {
		return err
	}
	roleM.MenusM = menus

	// 更新
	if err := b.ds.Roles().Update(ctx, roleM); err != nil {
		return err
	}

	// 更新需要重新绑定角色的权限
	err = b.bindPermissions(ctx, id, menus)
	if err != nil {
		return err
	}

	return nil
}

// Create 创建角色
func (b *roleBiz) Create(ctx context.Context, req *v1.CreateRoleRequest) error {
	var roleM model.RoleM
	_ = copier.Copy(&roleM, req)

	// 查询是否有该角色
	err := b.ds.DB().Where("id = ?", req.ID).First(&model.RoleM{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ErrRoleAlreadyExist
		}
	}

	// 添加角色与菜单的关系
	var menus []*model.MenuM
	if err := b.ds.DB().Preload("SysApisM").Where("id in (?)", req.MenuIds).Find(&menus).Error; err != nil {
		return err
	}
	roleM.MenusM = menus

	if err := b.ds.Roles().Create(ctx, &roleM); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'role.PRIMARY'", err.Error()); match {
			return errno.ErrRoleAlreadyExist
		}
		return err
	}

	// 绑定角色权限
	err = b.bindPermissions(ctx, roleM.ID, roleM.MenusM)
	if err != nil {
		return err
	}

	return nil
}

func (b *roleBiz) bindPermissions(ctx context.Context, roleId string, menus []*model.MenuM) error {
	// 给角色赋予权限
	var permissionInfos []map[string]string
	for _, item := range menus {
		menu := item
		for _, api := range menu.SysApisM {
			permissionInfos = append(permissionInfos, map[string]string{"Method": api.Method, "Path": api.Path})
		}
	}

	if len(permissionInfos) > 0 {
		// 添加角色权限
		err := b.ds.Casbins().CreateOrUpdate(ctx, roleId, permissionInfos)
		if err != nil {
			return err
		}
	}

	return nil
}

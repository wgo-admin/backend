package user

import (
	"context"
	"errors"
	"regexp"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/wgo-admin/backend/internal/app/store"
	"github.com/wgo-admin/backend/internal/app/store/model"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/known"
	"github.com/wgo-admin/backend/internal/pkg/log"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
	"github.com/wgo-admin/backend/pkg/auth"
	"github.com/wgo-admin/backend/pkg/token"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type IUserBiz interface {
	Register(ctx context.Context, req *v1.RegisterUserRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error)
	Logout(ctx context.Context, username string) error
	List(ctx context.Context, req *v1.UserListRequest) (*v1.UserListResponse, error)
	Get(ctx context.Context, username string) (*v1.UserInfo, error)
	Delete(ctx context.Context, username string) error
	Update(ctx context.Context, username string, req *v1.UpdateUserRequest) error
	ChangePassword(ctx context.Context, username string, req *v1.ChangePasswordRequest) error
}

var _ IUserBiz = (*userBiz)(nil)

func NewBiz(ds store.IStore) *userBiz {
	return &userBiz{ds}
}

type userBiz struct {
	ds store.IStore
}

// ChangePassword implements IUserBiz
func (b *userBiz) ChangePassword(ctx context.Context, username string, req *v1.ChangePasswordRequest) error {
	userM, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		return err
	}

	if err := auth.Compare(userM.Password, req.OldPassword); err != nil {
		return errno.ErrPasswordIncorrect
	}

	userM.Password, _ = auth.Encrypt(req.NewPassword)
	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return err
	}

	return nil
}

func (b *userBiz) Update(ctx context.Context, username string, req *v1.UpdateUserRequest) error {
	userM, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		return err
	}

	if req.Email != nil {
		userM.Email = *req.Email
	}

	if req.Nickname != nil {
		userM.Nickname = *req.Nickname
	}

	if req.Phone != nil {
		userM.Phone = *req.Phone
	}

	if req.RoleID != nil {
		_, err := b.ds.Roles().Get(ctx, *req.RoleID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errno.ErrRoleNotFound
			}
			return err
		}
		userM.RoleID = *req.RoleID
	}

	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return err
	}

	return nil
}

// Delete 是 UserBiz 接口中 `Delete` 方法的实现.
func (b *userBiz) Delete(ctx context.Context, username string) error {
	if err := b.ds.Users().Delete(ctx, username); err != nil {
		return err
	}

	return nil
}

// 获取用户详情
func (b *userBiz) Get(ctx context.Context, username string) (*v1.UserInfo, error) {
	user, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}

		return nil, err
	}

	var resp v1.UserInfo
	copier.Copy(&resp, user)

	resp.CreatedAt = user.CreatedAt.Format("2006-01-02 15:04:05")
	resp.UpdatedAt = user.UpdatedAt.Format("2006-01-02 15:04:05")

	return &resp, nil
}

// List 查询用户列表
func (b *userBiz) List(ctx context.Context, req *v1.UserListRequest) (*v1.UserListResponse, error) {
	var userM model.UserM
	copier.Copy(&userM, req)
	total, list, err := b.ds.Users().List(ctx, &userM, req.Offset(), req.Limit())
	if err != nil {
		log.C(ctx).Errorw("Failed to list users from storage", "err", err)
		return nil, err
	}

	var m sync.Map
	eg, ctx := errgroup.WithContext(ctx)
	// 使用 goroutine 提高接口性能
	for _, item := range list {
		user := item
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				role, err := b.ds.Roles().Get(ctx, user.RoleID)
				if err != nil {
					log.C(ctx).Errorw("Failed to role detail", "err", err)
					return err
				}

				m.Store(user.ID, &v1.UserInfo{
					Username:  user.Username,
					Nickname:  user.Nickname,
					Email:     user.Email,
					Phone:     user.Phone,
					RoleName:  role.Name,
					RoleID:    user.RoleID,
					CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
					UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
				})

				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.C(ctx).Errorw("Failed to wait all function calls returned", "err", err)
		return nil, err
	}

	users := make([]*v1.UserInfo, 0, len(list))
	for _, item := range list {
		user, _ := m.Load(item.ID)
		users = append(users, user.(*v1.UserInfo))
	}

	log.C(ctx).Debugw("Get users from backend storage", "count", len(users))

	return &v1.UserListResponse{Total: total, List: users}, nil
}

// Logout 用户退出业务逻辑
func (b *userBiz) Logout(ctx context.Context, username string) error {
	if err := b.ds.Cache().Del(ctx, known.GetRedisUserKey(username)); err != nil {
		return err
	}
	return nil
}

// Register 注册用户的业务逻辑
func (b *userBiz) Register(ctx context.Context, req *v1.RegisterUserRequest) error {
	var userM model.UserM
	_ = copier.Copy(&userM, req)

	var err error
	// 开启一个事务
	tx := b.ds.DB().Begin()
	// 有错误则回滚，否则就提交
	defer func() {
		if err != nil {
			if err = tx.Rollback().Error; err != nil {
				log.Errorw("rollback error", "error", err)
			}
		} else {
			if err := tx.Commit().Error; err != nil {
				log.Errorw("commit error", "error", err)
			}
		}
	}()

	// 查询是否有该角色
	var roleM model.RoleM
	if err = tx.First(&roleM, "id = ?", req.RoleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ErrRoleNotFound
		}
		return err
	}

	// 创建一个 user
	if err = tx.Create(&userM).Error; err != nil {
		log.Errorw("create user error", "error", err)
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'user.username'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}
		return err
	}

	return nil
}

// Login 用户登录业务逻辑
func (b *userBiz) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	// 1.获取登录用户的信息
	user, err := b.ds.Users().Get(ctx, req.Username)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	// 2.对比传入的明文密码和数据库中的加密密码
	if err := auth.Compare(user.Password, req.Password); err != nil {
		return nil, errno.ErrPasswordIncorrect
	}

	// 3.如果密码校验成功则签发token
	t, err := token.Sign(user.Username, user.RoleID)
	if err != nil {
		return nil, errno.ErrSignToken
	}

	// 登录的用户缓存到redis
	err = b.ds.Cache().SetStruct(ctx, known.GetRedisUserKey(user.Username), user, 0)
	if err != nil {
		return nil, err
	}

	return &v1.LoginResponse{Token: t}, nil
}

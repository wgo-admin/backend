package user

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/wgo-admin/backend/internal/app/store"
	"github.com/wgo-admin/backend/internal/app/store/model"
	"github.com/wgo-admin/backend/internal/pkg/errno"
	"github.com/wgo-admin/backend/internal/pkg/known"
	"github.com/wgo-admin/backend/internal/pkg/log"
	v1 "github.com/wgo-admin/backend/pkg/api/app/v1"
	"github.com/wgo-admin/backend/pkg/auth"
	"github.com/wgo-admin/backend/pkg/token"
	"gorm.io/gorm"
	"regexp"
)

type IUserBiz interface {
	Register(ctx context.Context, req *v1.RegisterUserRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error)
	Logout(ctx context.Context, username string) error
}

var _ IUserBiz = (*userBiz)(nil)

func NewBiz(ds store.IStore) *userBiz {
	return &userBiz{ds}
}

type userBiz struct {
	ds store.IStore
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

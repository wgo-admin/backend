package v1

// RegisterUserRequest 定义 `POST /v1/users/register` 接口的请求参数
type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=6,max=20"`
	Nickname string `json:"nickname" validate:"required,min=1,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,min=11,max=11"`
	RoleID   string `json:"roleId" validate:"required"`
}

// LoginRequest 定义 `POST /v1/users/login` 接口的请求参数
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse 定义 `POST /v1/users/login` 接口的返回参数
type LoginResponse struct {
	Token string `json:"token"`
}

// 定义用户的详细信息
type UserInfo struct {
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	RoleID    string `json:"roleId"`
	RoleName  string `json:"roleName"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// 定义 `GET /v1/users` 接口的请求参数
type UserListRequest struct {
	Pager
	Username string `form:"username"`
	Email    string `form:"email"`
}

// 定义 `GET /v1/users` 接口的返回参数
type UserListResponse struct {
	Total int64       `json:"total"`
	List  []*UserInfo `json:"list"`
}

// 定义 `PUT /v1/users` 接口请求参数
type UpdateUserRequest struct {
	Nickname *string `json:"nickname" validate:"required,min=1,max=255"`
	Email    *string `json:"email" validate:"required,email"`
	Phone    *string `json:"phone" validate:"required,min=11,max=11"`
	RoleID   *string `json:"roleId" validate:"required"`
}

// ChangePasswordRequest 指定了 `PUT /v1/users/{name}/change-password` 接口的请求参数.
type ChangePasswordRequest struct {
	// 旧密码.
	OldPassword string `json:"oldPassword" validate:"required"`

	// 新密码.
	NewPassword string `json:"newPassword" validate:"required,min=6,max=20"`
}

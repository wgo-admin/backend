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

package v1

// CreateRoleRequest 定义创建角色 `POST /v1/roles` 的请求参数
type CreateRoleRequest struct {
	ID      string  `json:"id" validate:"required"`
	Name    string  `json:"name" validate:"required"`
	Desc    string  `json:"desc" validate:"required"`
	MenuIds []int64 `json:"menuIds" validate:"required"`
}

// UpdateRoleRequest 定义更新角色 `PUT /v1/roles/:id` 的请求参数
type UpdateRoleRequest struct {
	Name    string  `json:"name" validate:"required"`
	Desc    string  `json:"desc" validate:"required"`
	MenuIds []int64 `json:"menuIds" validate:"required"`
}

type RoleInfo struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	Name      string  `json:"name"`
	Desc      string  `json:"desc"`
	UserCount int     `json:"userCount"` // 绑定该角色的用户数量
	MenuIds   []int64 `json:"menuIds"`   // 关联的菜单id切片
}

// QueryRoleRequest 定义查询角色列表 `GET /v1/roles` 的请求参数
type QueryRoleRequest struct {
	Pager
	ID   string `form:"id"`
	Name string `form:"name"`
}

// QueryRoleResponse 定义查询角色列表 `GET /v1/roles` 的响应结果
type QueryRoleResponse struct {
	Total int64       `json:"total"`
	List  []*RoleInfo `json:"list"`
}

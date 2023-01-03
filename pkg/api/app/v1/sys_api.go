package v1

// 定义 `POST /v1/sys_apis` 的请求参数
type CreateSysApiRequest struct {
	Title     string `json:"title" validate:"required"`
	Method    string `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Path      string `json:"path" validate:"required"`
	GroupName string `json:"groupName" validate:"required"`
}

type SysApiInfo struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Title     string `json:"title"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	GroupName string `json:"groupName"`
}

// 定义 `GET /v1/sys_apis` 的请求参数
type QuerySysApiRequest struct {
	Pager
	Method    string `form:"method" json:"method" validate:"omitempty,oneof=GET POST PUT PATCH DELETE"`
	Path      string `form:"path" json:"path"`
	GroupName string `form:"groupName" json:"groupName"`
}

// 定义 `GET /v1/sys_apis` 的返回参数
type QuerySysApiResponse struct {
	Total int64         `json:"total"`
	List  []*SysApiInfo `json:"list"`
}

// 定义 `PUT /v1/sys_apis/:id` 的请求参数
type UpdateSysApiRequest struct {
	Title     *string `json:"title" validate:"required"`
	Method    *string `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Path      *string `json:"path" validate:"required"`
	GroupName *string `json:"groupName" validate:"required"`
}

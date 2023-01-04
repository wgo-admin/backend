package v1

// 定义创建或更新菜单 `POST /v1/menus` `PUT /v1/menus/:id` 的请求参数
type CreateOrUpdateMenuRequest struct {
	ParentID  *int64  `json:"parentId" validate:"omitempty"`
	Title     string  `json:"title" validate:"required"`
	Sort      *int    `json:"sort" validate:"required"`
	Type      string  `json:"type" validate:"required,oneof=C M B"`
	Icon      string  `json:"icon" validate:"omitempty"`
	Component string  `json:"component" validate:"omitempty"`
	Path      string  `json:"path" validate:"omitempty"`
	Perm      string  `json:"perm" validate:"omitempty"`
	Hidden    bool    `json:"hidden" validate:"omitempty"`
	IsLink    bool    `json:"isLink" validate:"omitempty"`
	KeepAlive bool    `json:"keepAlive" validate:"omitempty"`
	ApiIds    []int64 `json:"apiIds" validate:"omitempty"`
}

type MenuInfo struct {
	ID        int64       `json:"id"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
	ParentID  int64       `json:"parentId"`
	Title     string      `json:"title"`
	Sort      int         `json:"sort"`
	Type      string      `json:"type"`
	Icon      string      `json:"icon,"`
	Component string      `json:"component"`
	Path      string      `json:"path"`
	Perm      string      `json:"perm"`
	Hidden    bool        `json:"hidden"`
	IsLink    bool        `json:"isLink"`
	KeepAlive bool        `json:"keepAlive"`
	ApiIds    []int64     `json:"apiIds"`
	Children  []*MenuInfo `json:"chidlren,omitempty"`
	IsLeaf    bool        `json:"isLeaf,omitempty"`
}

// 定义菜单树 `GET /v1/menus/tree` 响应参数
type QueryMenuTreeResponse struct {
	List []*MenuInfo `json:"list"`
}

// 定义获取用户的菜单 `GET /v1/menus/user` 响应参数
type QueryUserMenuResponse struct {
	List []*MenuInfo `json:"list"`
}

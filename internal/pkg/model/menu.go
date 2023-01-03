package model

type MenuM struct {
	BaseM
	Title     string    `gorm:"size:128;comment:标题"`
	Sort      int       `gorm:"size:4;comment:排序"`
	Type      string    `gorm:"size:1;comment:菜单类型, C:代表目录 M:代表菜单 B:代表按钮"`
	ParentID  int64     `gorm:"comment:父级ID"`
	Icon      string    `gorm:"size:128;comment:图标"`
	Component string    `gorm:"size:255;comment:组件路径"`
	Path      string    `gorm:"size:128;comment:路由路径"`
	Perm      string    `gorm:"size:255;comment:权限标识"`
	Hidden    bool      `gorm:"comment:是否隐藏"`
	IsLink    bool      `gorm:"comment:是否是链接"`
	KeepAlive bool      `gorm:"comment:是否缓存"`
	SysApisM  []SysApiM `gorm:"many2many:menu_sys_apis;foreignKey:ID;joinForeignKey:menu_id;references:ID;joinReferences:sys_api_id;"`
}

func (a *MenuM) TableName() string {
	return "menu"
}

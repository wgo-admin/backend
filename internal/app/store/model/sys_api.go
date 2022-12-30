package model

type SysApiM struct {
	BaseM
	Title     string `gorm:"size:128;comment:标题"`
	Method    string `gorm:"size:16;index;comment:请求方法"`
	Path      string `gorm:"size:128;index;comment:请求资源路径"`
	GroupName string `gorm:"size:20;index;comment:组名"`
}

func (a *SysApiM) TableName() string {
	return "sys_api"
}

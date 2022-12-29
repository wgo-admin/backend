package v1

type PermissionInfo struct {
	Path   string `json:"path"`   // 请求路径
	Method string `json:"method"` // 请求方法
}

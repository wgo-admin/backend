package errno

import "net/http"

var (
	// ErrSysApiAlreadyExist 表示已存在接口.
	ErrSysApiAlreadyExist = &Errno{
		HTTP:    http.StatusBadRequest,
		Code:    "FailedOperation.SysApiAlreadyExist",
		Message: "SysApi already exist.",
	}
)

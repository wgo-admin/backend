package errno

var (
	// ErrRoleNotFound 表示未找到用户.
	ErrRoleNotFound = &Errno{
		HTTP:    404,
		Code:    "ResourceNotFound.RoleNotFound",
		Message: "Role was not found.",
	}

	ErrRoleAlreadyExist = &Errno{
		HTTP:    404,
		Code:    "FailedOperation.RoleAlreadyExist",
		Message: "Role already exist.",
	}
)

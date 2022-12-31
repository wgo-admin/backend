package errno

import "net/http"

var (
	// ErrUserAlreadyExist 代表用户已存在
	ErrUserAlreadyExist = &Errno{
		HTTP:    http.StatusBadRequest,
		Code:    "FailedOperation.UserAlreadyExist",
		Message: "User already exist.",
	}

	// ErrUserNotFound 表示未找到用户.
	ErrUserNotFound = &Errno{
		HTTP:    404,
		Code:    "ResourceNotFound.UserNotFound",
		Message: "User was not found.",
	}

	// ErrPasswordIncorrect 表示密码不正确.
	ErrPasswordIncorrect = &Errno{
		HTTP:    401,
		Code:    "InvalidParameter.PasswordIncorrect",
		Message: "Password was incorrect.",
	}

	ErrUserLogout = &Errno{
		HTTP:    401,
		Code:    "FailedOperation.UserLogout",
		Message: "User is logout.",
	}
)

package errno

import "net/http"

var (
  ErrHasSameApi = &Errno{
    HTTP:    http.StatusBadRequest,
    Code:    "FailedOperation.HasSameApi",
    Message: "Two identical api's exist",
  }

  ErrClearPermissionFailed = &Errno{
    HTTP:    http.StatusBadRequest,
    Code:    "FailedOperation.ClearPermissionFailed",
    Message: "Clear permission failed.",
  }
)

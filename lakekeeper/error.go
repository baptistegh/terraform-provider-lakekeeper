package lakekeeper

import (
	"net/http"

	"github.com/hashicorp/errwrap"
)

type ApiError struct {
	Code    int
	Message string
}

func (e *ApiError) Error() string {
	return e.Message
}

func ErrorIs404(err error) bool {
	lakekeeperError, ok := errwrap.GetType(err, &ApiError{}).(*ApiError)

	return ok && lakekeeperError != nil && lakekeeperError.Code == http.StatusNotFound
}

func ErrorIs409(err error) bool {
	lakekeeperError, ok := errwrap.GetType(err, &ApiError{}).(*ApiError)

	return ok && lakekeeperError != nil && lakekeeperError.Code == http.StatusConflict
}

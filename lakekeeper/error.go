package lakekeeper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ApiError struct {
	Status     string         `json:"-"`
	StatusCode int            `json:"-"`
	Message    string         `json:"-"`
	Response   *ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Stack   []string `json:"stack"`
	Type    string   `json:"type"`
}

func (e *ApiError) Error() string {
	if e.Response == nil {
		return fmt.Sprintf("unexpected error response: %s", e.Message)
	}
	return fmt.Sprintf("api error, code=%d message=%s type=%s", e.Response.Code, e.Response.Message, e.Response.Type)
}

func (e *ApiError) IsAuthError() bool {
	return e.StatusCode == http.StatusUnauthorized || e.StatusCode == http.StatusForbidden
}

func ApiErrorFromResponse(response *http.Response) *ApiError {
	defer response.Body.Close()

	apiErr := ApiError{}
	apiErr.Status = response.Status
	apiErr.StatusCode = response.StatusCode

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return &apiErr
	}

	if len(responseBody) > 0 {
		_ = json.Unmarshal(responseBody, &apiErr)
	}

	return &apiErr
}

func ApiErrorFromError(format string, a ...any) *ApiError {
	return &ApiError{
		Message: fmt.Sprintf(format, a...),
	}
}

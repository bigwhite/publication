package errs

import "net/http"

// ErrorType 定义业务错误码类型
type ErrorType string

const (
	// 通用错误
	ErrInternalServer ErrorType = "INTERNAL_SERVER_ERROR"
	ErrBadRequest     ErrorType = "BAD_REQUEST"
	ErrNotFound       ErrorType = "NOT_FOUND"
	ErrUnauthorized   ErrorType = "UNAUTHORIZED"

	// 业务特定错误 (Example)
	ErrUserFrozen    ErrorType = "USER_FROZEN"
	ErrQuotaExceeded ErrorType = "QUOTA_EXCEEDED"
)

// Map ErrorType to HTTP Status Code
var statusMap = map[ErrorType]int{
	ErrInternalServer: http.StatusInternalServerError,
	ErrBadRequest:     http.StatusBadRequest,
	ErrNotFound:       http.StatusNotFound,
	ErrUnauthorized:   http.StatusUnauthorized,
	ErrUserFrozen:     http.StatusForbidden,
	ErrQuotaExceeded:  http.StatusTooManyRequests,
}

func (t ErrorType) HTTPStatus() int {
	if code, ok := statusMap[t]; ok {
		return code
	}
	return http.StatusInternalServerError
}

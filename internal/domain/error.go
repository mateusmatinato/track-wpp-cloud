package domain

import (
	"fmt"
	"net/http"
	"strings"
)

type ApiError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func NewInternalServerError(message string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusInternalServerError,
		Code:       "internal server err",
		Message:    message,
	}
}

func NewUnauthorizedError(message string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusUnauthorized,
		Code:       "unauthorized",
		Message:    message,
	}
}

func (err *ApiError) Error() string {
	code := strings.ReplaceAll(strings.ToLower(http.StatusText(err.StatusCode)), " ", "_")
	return fmt.Sprintf("%d %s: %s", err.StatusCode, code, err.Message)
}

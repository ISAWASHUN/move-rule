package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

type ErrorResponse struct {
	Error AppError `json:"error"`
}

func NewBadRequestError(msg string, err error) *AppError {
	return &AppError{
		Code:       "bad_request",
		Message:    msg,
		StatusCode: http.StatusBadRequest,
		Err:        err,
	}
}

func NewNotFoundError(msg string, err error) *AppError {
	return &AppError{
		Code:       "not_found",
		Message:    msg,
		StatusCode: http.StatusNotFound,
		Err:        err,
	}
}

func NewInternalError(msg string, err error) *AppError {
	return &AppError{
		Code:       "internal_error",
		Message:    msg,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func HandleError(c *gin.Context, err *AppError) {
	c.AbortWithStatusJSON(err.StatusCode, ErrorResponse{
		Error: *err,
	})
}

func JSON(c *gin.Context, status int, err *AppError) {
	c.AbortWithStatusJSON(status, gin.H{
		"error": err,
	})
}

func JSONInternal(c *gin.Context, msg string) {
	JSON(c, http.StatusInternalServerError, NewInternalError(msg, nil))
}

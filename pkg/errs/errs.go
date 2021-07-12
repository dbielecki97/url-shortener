package errs

import (
	"fmt"
	"net/http"
)

type RestErr interface {
	Code() int
	Message() string
	Err() string
	Error() string
}

type Err struct {
	StatusCode int    `json:"code,omitempty"`
	Msg        string `json:"message,omitempty"`
	ErrMessage string `json:"error,omitempty"`
}

func (e Err) Code() int {
	return e.StatusCode
}

func (e Err) Message() string {
	return e.Msg
}

func (e Err) Err() string {
	return e.ErrMessage
}

func (e Err) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s", e.Message(), e.Code(), e.ErrMessage)
}

func (e Err) AsMessage() *Err {
	return &Err{Msg: e.Msg}
}

func NewErr(msg string, code int, err string) RestErr {
	return &Err{
		Msg:        msg,
		StatusCode: code,
		ErrMessage: err,
	}
}

func NewUnexpectedError(message string) *Err {
	return &Err{
		StatusCode: http.StatusInternalServerError,
		Msg:        message,
	}
}

func NewValidationError(message string) *Err {
	return &Err{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        message,
	}
}
func NewNotFoundError(message string) *Err {
	return &Err{
		StatusCode: http.StatusNotFound,
		Msg:        message,
	}
}

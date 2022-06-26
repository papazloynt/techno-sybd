package constants

import (
	"errors"
	"net/http"
)

// CodedError is an error wrapper which wraps errors with http status codes.
type CodedError struct {
	err  error
	code int
}

func (ce *CodedError) Error() string {
	return ce.err.Error()
}

func (ce *CodedError) Code() int {
	return ce.code
}

func CreateNewError(s string, code int) *CodedError {
	return &CodedError{errors.New(s), code}
}

var (
	// Bad Request
	ErrBindRequest     = &CodedError{errors.New("failed to bind request"), http.StatusBadRequest}
	ErrValidateRequest = &CodedError{errors.New("failed to validate request"), http.StatusBadRequest}
	ErrDBNotFound      = &CodedError{errors.New("not found in the database"), http.StatusBadRequest}
	ErrBadJson         = &CodedError{errors.New("bad json request"), http.StatusBadRequest}

	// User
	ErrUserAlreadyExists = &CodedError{errors.New("user with this nickname or email already exists"), http.StatusConflict}
)

var (
	ParseError = map[string]*CodedError{
		ErrBindRequest.Error():       ErrBindRequest,
		ErrValidateRequest.Error():   ErrValidateRequest,
		ErrDBNotFound.Error():        ErrDBNotFound,
		ErrBadJson.Error():           ErrBadJson,
		ErrUserAlreadyExists.Error(): ErrUserAlreadyExists,
	}
)

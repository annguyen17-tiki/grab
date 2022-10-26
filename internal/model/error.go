package model

import (
	"errors"
	"fmt"
)

var (
	ErrorCannotUnmarshalJSONB = errors.New("cannot_unmarshal_JSONB_value")
)

type ErrBadRequest struct {
	msg string
}

func (e *ErrBadRequest) Error() string {
	return e.msg
}

func NewErrBadRequest(format string, args ...interface{}) error {
	return &ErrBadRequest{
		msg: fmt.Sprintf(format, args...),
	}
}

type ErrUnauthorized struct {
	msg string
}

func (e *ErrUnauthorized) Error() string {
	return e.msg
}

func NewErrUnauthorized(format string, args ...interface{}) error {
	return &ErrUnauthorized{
		msg: fmt.Sprintf(format, args...),
	}
}

type ErrForbidden struct {
	msg string
}

func (e *ErrForbidden) Error() string {
	return e.msg
}

func NewErrForbidden(format string, args ...interface{}) error {
	return &ErrForbidden{
		msg: fmt.Sprintf(format, args...),
	}
}

type ErrNotFound struct {
	msg string
}

func (e *ErrNotFound) Error() string {
	return e.msg
}

func NewErrNotFound(format string, args ...interface{}) error {
	return &ErrNotFound{
		msg: fmt.Sprintf(format, args...),
	}
}

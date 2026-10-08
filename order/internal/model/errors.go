package model

import "errors"

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrCannotProcessOrder  = errors.New("cannot process order")
	ErrInternalServerError = errors.New("internal server error")
)

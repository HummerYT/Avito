package models

import "errors"

var (
	ErrAuthUser   = errors.New("user is not authorized")
	ErrValidation = errors.New("validation error")
)

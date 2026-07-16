package service

import "errors"

var (
	ErrUnauthorized = errors.New("only supervisors can perform this action")
)

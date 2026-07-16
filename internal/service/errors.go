package service

import "errors"

var (
	ErrUnauthorized      = errors.New("only supervisors can perform this action")
	ErrWorkerOnly        = errors.New("only workers can perform this action")
	ErrDifferentWorker   = errors.New("only the worker of this task can perform this action")
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidTransition = errors.New("task status transition is invalid")
)

package app

import "errors"

var (
	ErrObjectNotFound = errors.New("object not found")
	ErrUserNotFound   = errors.New("user not found")
	ErrCreateReview   = errors.New("failed to create review")
)

package business

import "errors"

var (
	ErrInvalidName  = errors.New("business: invalid name")
	ErrInvalidToken = errors.New("business: invalid token")
	ErrNotFound     = errors.New("business: not found")
)

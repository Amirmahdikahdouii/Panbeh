package otp

import "errors"

var (
	ErrInvalidPhone    = errors.New("otp: invalid phone number")
	ErrInvalidTTL      = errors.New("otp: invalid ttl")
	ErrInvalidCode     = errors.New("otp: invalid code")
	ErrInvalidBusiness = errors.New("otp: invalid business id")
	ErrNotFound        = errors.New("otp: not found")
)

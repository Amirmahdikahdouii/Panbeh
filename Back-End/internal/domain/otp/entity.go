package otp

import (
	"regexp"
	"time"
)

type IranPhoneNumber string

var iranPhoneRegex = regexp.MustCompile(`^(?:\+98|0)?9\d{9}$`)

func NewIranPhoneNumber(value string) (IranPhoneNumber, error) {
	if !iranPhoneRegex.MatchString(value) {
		return "", ErrInvalidPhone
	}
	return IranPhoneNumber(value), nil
}

type CodeTTL time.Duration

func NewCodeTTL(ttl time.Duration) (CodeTTL, error) {
	if ttl <= 0 {
		return CodeTTL(ttl), ErrInvalidTTL
	}
	return CodeTTL(ttl), nil
}

type OTP struct {
	BusinessID  string
	PhoneNumber IranPhoneNumber
	Code        string
	ExpiresAt   time.Time
}

func (o OTP) Expired(now time.Time) bool {
	return !o.ExpiresAt.After(now)
}

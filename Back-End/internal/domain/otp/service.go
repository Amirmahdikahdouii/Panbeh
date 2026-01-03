package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"
)

var (
	codeRe = regexp.MustCompile(`^\d{6}$`)
)

type Service struct {
	now     func() time.Time
	ttl     CodeTTL
	codeGen func() (string, error)
}

type ServiceConfig struct {
	Now     func() time.Time
	TTL     CodeTTL
	CodeGen func() (string, error)
}

func NewService(cfg ServiceConfig) *Service {
	now := cfg.Now
	if now == nil {
		now = time.Now
	}
	codeGen := cfg.CodeGen
	if codeGen == nil {
		codeGen = defaultCode
	}
	return &Service{
		now:     now,
		ttl:     cfg.TTL,
		codeGen: codeGen,
	}
}

func (s *Service) NewOTP(businessID string, phone IranPhoneNumber) (OTP, error) {
	if strings.TrimSpace(businessID) == "" {
		return OTP{}, ErrInvalidBusiness
	}
	code, err := s.codeGen()
	if err != nil {
		return OTP{}, err
	}
	if err := ValidateCode(code); err != nil {
		return OTP{}, err
	}

	now := s.now()
	return OTP{
		BusinessID:  businessID,
		PhoneNumber: phone,
		Code:        code,
		ExpiresAt:   now.Add(time.Duration(s.ttl)),
	}, nil
}

func (s *Service) Verify(stored OTP, businessID string, phone IranPhoneNumber, code string) (bool, error) {
	if err := ValidateCode(code); err != nil {
		return false, err
	}
	if strings.TrimSpace(businessID) == "" {
		return false, ErrInvalidBusiness
	}

	if stored.BusinessID != businessID || stored.PhoneNumber != phone {
		return false, nil
	}
	if stored.Expired(s.now()) {
		return false, nil
	}
	return stored.Code == code, nil
}

func ValidateCode(code string) error {
	code = strings.TrimSpace(code)
	if !codeRe.MatchString(code) {
		return ErrInvalidCode
	}
	return nil
}

func defaultCode() (string, error) {
	// 6-digit numeric OTP.
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

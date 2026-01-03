package business

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
)

type Service struct {
	now      func() time.Time
	idGen    func() (string, error)
	tokenGen func() (string, error)
}

type ServiceConfig struct {
	Now      func() time.Time
	IDGen    func() (string, error)
	TokenGen func() (string, error)
}

func NewService(cfg ServiceConfig) *Service {
	now := cfg.Now
	if now == nil {
		now = time.Now
	}
	idGen := cfg.IDGen
	if idGen == nil {
		idGen = defaultID
	}
	tokenGen := cfg.TokenGen
	if tokenGen == nil {
		tokenGen = defaultToken
	}
	return &Service{now: now, idGen: idGen, tokenGen: tokenGen}
}

func (s *Service) NewBusiness(name string) (Business, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Business{}, ErrInvalidName
	}

	id, err := s.idGen()
	if err != nil {
		return Business{}, err
	}

	token, err := s.tokenGen()
	if err != nil {
		return Business{}, err
	}

	return Business{
		ID:        id,
		Name:      name,
		Token:     token,
		CreatedAt: s.now(),
	}, nil
}

func ValidateToken(token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return ErrInvalidToken
	}
	// Token format is intentionally opaque; minimal validation prevents accidental empty/whitespace tokens.
	return nil
}

func defaultID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func defaultToken() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

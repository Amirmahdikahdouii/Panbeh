package business_test

import (
	"testing"
	"time"

	"github.com/panbeh/otp-backend/internal/domain/business"
)

func TestService_NewBusiness_ValidatesName(t *testing.T) {
	svc := business.NewService(business.ServiceConfig{
		Now: func() time.Time { return time.Unix(10, 0) },
		IDGen: func() (string, error) {
			return "id1", nil
		},
		TokenGen: func() (string, error) {
			return "tok1", nil
		},
	})

	if _, err := svc.NewBusiness("  "); err != business.ErrInvalidName {
		t.Fatalf("expected ErrInvalidName, got %v", err)
	}
}

func TestService_NewBusiness_CreatesDeterministicBusiness(t *testing.T) {
	now := time.Unix(10, 0)
	svc := business.NewService(business.ServiceConfig{
		Now: func() time.Time { return now },
		IDGen: func() (string, error) {
			return "id1", nil
		},
		TokenGen: func() (string, error) {
			return "tok1", nil
		},
	})

	b, err := svc.NewBusiness("Acme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.ID != "id1" || b.Token != "tok1" || b.Name != "Acme" || !b.CreatedAt.Equal(now) {
		t.Fatalf("unexpected business: %#v", b)
	}
}

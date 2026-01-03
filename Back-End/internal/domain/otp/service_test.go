package otp_test

import (
	"testing"
	"time"

	"github.com/panbeh/otp-backend/internal/domain/otp"
)

func Test_NewValidTTL(t *testing.T) {
	ttl, err := otp.NewCodeTTL(5 * time.Minute)
	if err != nil {
		t.Fatalf("expected No Error for create valid ttl, got %v", err)
	}
	if time.Duration(ttl) != 5*time.Minute {
		t.Fatalf("expected 5 minute duration ttl, got %v", ttl)
	}
}

func Test_NewInValidTTL(t *testing.T) {
	_, err := otp.NewCodeTTL(-5 * time.Minute)
	if err != otp.ErrInvalidTTL {
		t.Fatalf("expected ErrInvalidTTL Error for create invalid ttl, got %v", err)
	}
}

func TestService_NewOTP_ValidatesPhoneAndBusiness(t *testing.T) {
	ttl, err := otp.NewCodeTTL(5 * time.Minute)
	if err != nil {
		t.Fatalf("expected No Error for create valid ttl, got %v", err)
	}

	svc := otp.NewService(otp.ServiceConfig{
		TTL: ttl,
	})
	// TODO: remove check business and replace it with businessID type
	if _, err := svc.NewOTP("", "+15551234567"); err != otp.ErrInvalidBusiness {
		t.Fatalf("expected ErrInvalidBusiness, got %v", err)
	}
}

func TestService_NewOTP_GeneratesDeterministicOTP(t *testing.T) {
	ttl, err := otp.NewCodeTTL(2 * time.Minute)
	if err != nil {
		t.Fatalf("expected No Error for create valid ttl, got %v", err)
	}
	now := time.Unix(100, 0)
	svc := otp.NewService(otp.ServiceConfig{
		Now: func() time.Time { return now },
		TTL: ttl,
		CodeGen: func() (string, error) {
			return "123456", nil
		},
	})

	o, err := svc.NewOTP("b1", "+15551234567")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if o.Code != "123456" {
		t.Fatalf("expected code 123456, got %q", o.Code)
	}
	if !o.ExpiresAt.Equal(now.Add(2 * time.Minute)) {
		t.Fatalf("unexpected ExpiresAt: %v", o.ExpiresAt)
	}
}

func TestService_Verify(t *testing.T) {
	ttl, err := otp.NewCodeTTL(2 * time.Minute)
	if err != nil {
		t.Fatalf("expected No Error for create valid ttl, got %v", err)
	}
	now := time.Unix(100, 0)
	svc := otp.NewService(otp.ServiceConfig{
		Now: func() time.Time { return now },
		TTL: ttl,
	})

	stored := otp.OTP{
		BusinessID:  "b1",
		PhoneNumber: "+15551234567",
		Code:        "123456",
		ExpiresAt:   now.Add(1 * time.Minute),
	}

	ok, err := svc.Verify(stored, "b1", "+15551234567", "123456")
	if err != nil || !ok {
		t.Fatalf("expected ok, got ok=%v err=%v", ok, err)
	}

	ok, err = svc.Verify(stored, "b1", "+15551234567", "000000")
	if err != nil || ok {
		t.Fatalf("expected mismatch false,nil got ok=%v err=%v", ok, err)
	}

	expired := stored
	expired.ExpiresAt = now.Add(-time.Second)
	ok, err = svc.Verify(expired, "b1", "+15551234567", "123456")
	if err != nil || ok {
		t.Fatalf("expected expired false,nil got ok=%v err=%v", ok, err)
	}
}

func Test_NewIranPhoneNumber(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    otp.IranPhoneNumber
		expectError bool
	}{
		// Valid cases
		{
			name:        "valid 10-digit number starting with 9",
			input:       "9123456789",
			expected:    otp.IranPhoneNumber("9123456789"),
			expectError: false,
		},
		{
			name:        "valid 11-digit number with 0 prefix",
			input:       "09123456789",
			expected:    otp.IranPhoneNumber("09123456789"),
			expectError: false,
		},
		{
			name:        "valid 12-digit number with +98 prefix",
			input:       "+989123456789",
			expected:    otp.IranPhoneNumber("+989123456789"),
			expectError: false,
		},
		{
			name:        "valid another 10-digit number",
			input:       "9876543210",
			expected:    otp.IranPhoneNumber("9876543210"),
			expectError: false,
		},
		{
			name:        "valid minimum 10-digit number",
			input:       "9000000000",
			expected:    otp.IranPhoneNumber("9000000000"),
			expectError: false,
		},
		{
			name:        "valid maximum 10-digit number",
			input:       "9999999999",
			expected:    otp.IranPhoneNumber("9999999999"),
			expectError: false,
		},

		// Invalid cases
		{
			name:        "invalid - starts with 8 instead of 9",
			input:       "8123456789",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid - starts with 0 but wrong length",
			input:       "0912345678",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid - too few digits (9 digits)",
			input:       "912345678",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid - too many digits (11 digits)",
			input:       "91234567890",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid - contains non-numeric characters",
			input:       "912345678a",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid - contains spaces",
			input:       "912 345 6789",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid - wrong +98 prefix with insufficient digits",
			input:       "+98912345678",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid - wrong +98 prefix with extra digits",
			input:       "+9891234567890",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid - starts with +99 instead of +98",
			input:       "+999123456789",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid - empty string",
			input:       "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid - only prefix +98",
			input:       "+98",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid - only prefix 0",
			input:       "0",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid - starts with 09 but wrong total length",
			input:       "091234567890",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid - special characters",
			input:       "9-123-456-789",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid - mixed alphanumeric",
			input:       "9abc456789",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := otp.NewIranPhoneNumber(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for input %q, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error for input %q, but got: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("expected %q, but got %q", tt.expected, result)
				}
			}
		})
	}
}

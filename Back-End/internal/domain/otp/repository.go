package otp

import "context"

type Repository interface {
	Save(ctx context.Context, otp OTP) error
	Get(ctx context.Context, businessID, phone string) (OTP, error)
	Delete(ctx context.Context, businessID, phone string) error

	// Consume atomically verifies the provided code and deletes the OTP if it matches.
	// This is required to guarantee single-use semantics under concurrent verification attempts.
	Consume(ctx context.Context, businessID, phone, code string) (bool, error)
}

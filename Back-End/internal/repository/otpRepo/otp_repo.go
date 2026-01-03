package otp

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/panbeh/otp-backend/internal/domain/otp"
)

type OTPRepository struct {
	client redis.UniversalClient
}

func NewOTPRepository(client redis.UniversalClient) otp.Repository {
	return &OTPRepository{client: client}
}

type otpPayload struct {
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (r *OTPRepository) Save(ctx context.Context, o otp.OTP) error {
	key := otpKey(o.BusinessID, o.PhoneNumber)
	ttl := time.Until(o.ExpiresAt)
	if ttl <= 0 {
		return nil
	}
	b, err := json.Marshal(otpPayload{Code: o.Code, ExpiresAt: o.ExpiresAt})
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, b, ttl).Err()
}

func (r *OTPRepository) Get(ctx context.Context, businessID string, phone otp.IranPhoneNumber) (otp.OTP, error) {
	key := otpKey(businessID, phone)
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return otp.OTP{}, otp.ErrNotFound
		}
		return otp.OTP{}, err
	}

	var p otpPayload
	if err := json.Unmarshal(val, &p); err != nil {
		return otp.OTP{}, err
	}

	return otp.OTP{
		BusinessID:  businessID,
		PhoneNumber: phone,
		Code:        p.Code,
		ExpiresAt:   p.ExpiresAt,
	}, nil
}

func (r *OTPRepository) Delete(ctx context.Context, businessID string, phone otp.IranPhoneNumber) error {
	return r.client.Del(ctx, otpKey(businessID, phone)).Err()
}

func (r *OTPRepository) Consume(ctx context.Context, businessID string, phone otp.IranPhoneNumber, code string) (bool, error) {
	// Atomic compare-and-delete to guarantee single-use OTP.
	const script = `
local key = KEYS[1]
local expected = ARGV[1]
local val = redis.call("GET", key)
if not val then
  return 0
end
local decoded = cjson.decode(val)
if decoded["code"] ~= expected then
  return 0
end
redis.call("DEL", key)
return 1
`
	res, err := r.client.Eval(ctx, script, []string{otpKey(businessID, phone)}, code).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

func otpKey(businessID string, phone otp.IranPhoneNumber) string {
	return "otp:" + businessID + ":" + string(phone)
}

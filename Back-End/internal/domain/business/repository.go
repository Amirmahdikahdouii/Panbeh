package business

import "context"

type Repository interface {
	Create(ctx context.Context, b Business) (Business, error)
	GetByToken(ctx context.Context, token string) (Business, error)
}

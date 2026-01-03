package business

import (
	"context"
	"database/sql"

	"github.com/panbeh/otp-backend/internal/domain/business"
)

type BusinessRepository struct {
	db *sql.DB
}

func NewBusinessRepository(db *sql.DB) business.Repository {
	return &BusinessRepository{db: db}
}

func (r *BusinessRepository) Create(ctx context.Context, b business.Business) (business.Business, error) {
	// ID/CreatedAt are generated in the domain service; repository persists them as-is.
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO businesses (id, name, token, created_at)
		VALUES ($1, $2, $3, $4)
	`, b.ID, b.Name, b.Token, b.CreatedAt)
	if err != nil {
		return business.Business{}, err
	}
	return b, nil
}

func (r *BusinessRepository) GetByToken(ctx context.Context, token string) (business.Business, error) {
	var b business.Business
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, token, created_at
		FROM businesses
		WHERE token = $1
	`, token).Scan(&b.ID, &b.Name, &b.Token, &b.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return business.Business{}, business.ErrNotFound
		}
		return business.Business{}, err
	}
	return b, nil
}

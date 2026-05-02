package repository

import (
	"backend/internal/models"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ApiKeyRepository struct {
	db *pgxpool.Pool
}

func NewApiKeyRepository(db *pgxpool.Pool) *ApiKeyRepository {
	return &ApiKeyRepository{db: db}
}

func (r *ApiKeyRepository) Create(ctx context.Context, key *models.ApiKey) error {
	query := `
		INSERT INTO api_keys (id, user_id, name, public_key, hashed_secret, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(ctx, query,
		key.ID,
		key.UserID,
		key.Name,
		key.PublicKey,
		key.HashedSecret,
		key.CreatedAt,
	)
	return err
}

func (r *ApiKeyRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.ApiKey, error) {
	query := `
		SELECT id, user_id, name, public_key, created_at, revoked_at 
		FROM api_keys 
		WHERE user_id = $1 AND revoked_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []models.ApiKey
	for rows.Next() {
		var k models.ApiKey
		if err := rows.Scan(&k.ID, &k.UserID, &k.Name, &k.PublicKey, &k.CreatedAt, &k.RevokedAt); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

func (r *ApiKeyRepository) Revoke(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE api_keys SET revoked_at = NOW() WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, id, userID)
	return err
}

func (r *ApiKeyRepository) UpdateSecret(ctx context.Context, id uuid.UUID, userID uuid.UUID, newHashedSecret string) error {
	query := `UPDATE api_keys SET hashed_secret = $1 WHERE id = $2 AND user_id = $3 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, query, newHashedSecret, id, userID)
	return err
}

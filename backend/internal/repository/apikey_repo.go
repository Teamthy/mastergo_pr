package repository

import (
	"backend/internal/models"
	"context"
	"log"

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
	log.Printf("ApiKeyRepository.Create: inserting key %s for user %s", key.ID, key.UserID)

	query := `
		INSERT INTO api_keys (id, user_id, name, public_key, hashed_secret, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	result, err := r.db.Exec(ctx, query,
		key.ID,
		key.UserID,
		key.Name,
		key.PublicKey,
		key.HashedSecret,
		key.CreatedAt,
	)

	if err != nil {
		log.Printf("ApiKeyRepository.Create: ERROR inserting key: %v", err)
		return err
	}

	log.Printf("ApiKeyRepository.Create: successfully inserted key %s, rows affected: %s", key.ID, result.String())
	return nil
}

func (r *ApiKeyRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.ApiKey, error) {
	log.Printf("ApiKeyRepository.ListByUserID: fetching keys for user %s", userID)

	query := `
		SELECT id, user_id, name, public_key, created_at, revoked_at 
		FROM api_keys 
		WHERE user_id = $1 AND revoked_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		log.Printf("ApiKeyRepository.ListByUserID: ERROR querying database: %v", err)
		return nil, err
	}
	defer rows.Close()

	var keys []models.ApiKey
	for rows.Next() {
		var k models.ApiKey
		if err := rows.Scan(&k.ID, &k.UserID, &k.Name, &k.PublicKey, &k.CreatedAt, &k.RevokedAt); err != nil {
			log.Printf("ApiKeyRepository.ListByUserID: ERROR scanning row: %v", err)
			return nil, err
		}
		keys = append(keys, k)
	}

	log.Printf("ApiKeyRepository.ListByUserID: found %d keys for user %s", len(keys), userID)
	return keys, nil
}

func (r *ApiKeyRepository) Revoke(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	log.Printf("ApiKeyRepository.Revoke: revoking key %s for user %s", id, userID)

	query := `UPDATE api_keys SET revoked_at = NOW() WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)

	if err != nil {
		log.Printf("ApiKeyRepository.Revoke: ERROR revoking key: %v", err)
		return err
	}

	log.Printf("ApiKeyRepository.Revoke: successfully revoked key %s, rows affected: %s", id, result.String())
	return nil
}

func (r *ApiKeyRepository) UpdateSecret(ctx context.Context, id uuid.UUID, userID uuid.UUID, newHashedSecret string) error {
	log.Printf("ApiKeyRepository.UpdateSecret: updating secret for key %s for user %s", id, userID)

	query := `UPDATE api_keys SET hashed_secret = $1 WHERE id = $2 AND user_id = $3 AND revoked_at IS NULL`
	result, err := r.db.Exec(ctx, query, newHashedSecret, id, userID)

	if err != nil {
		log.Printf("ApiKeyRepository.UpdateSecret: ERROR updating secret: %v", err)
		return err
	}

	log.Printf("ApiKeyRepository.UpdateSecret: successfully updated secret for key %s, rows affected: %s", id, result.String())
	return nil
}

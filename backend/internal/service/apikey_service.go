package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ApiKeyService struct {
	repo *repository.ApiKeyRepository
}

func NewApiKeyService(repo *repository.ApiKeyRepository) *ApiKeyService {
	return &ApiKeyService{repo: repo}
}

func (s *ApiKeyService) CreateKey(ctx context.Context, userID uuid.UUID, name string) (*models.ApiKeyCreateResponse, error) {
	rawSecret := fmt.Sprintf("sk_live_%s", generateSecureRandom(32))
	pubKey := fmt.Sprintf("pk_live_%s", generateSecureRandom(16))

	hash := sha256.Sum256([]byte(rawSecret))
	hashedSecret := hex.EncodeToString(hash[:])

	apiKey := &models.ApiKey{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         name,
		PublicKey:    pubKey,
		HashedSecret: hashedSecret,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, apiKey); err != nil {
		return nil, err
	}

	return &models.ApiKeyCreateResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		PublicKey: apiKey.PublicKey,
		SecretKey: rawSecret,
		CreatedAt: apiKey.CreatedAt,
	}, nil
}

func (s *ApiKeyService) ListKeys(ctx context.Context, userID uuid.UUID) ([]models.ApiKey, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *ApiKeyService) RevokeKey(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.repo.Revoke(ctx, id, userID)
}

func (s *ApiKeyService) RegenerateKey(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.ApiKeyRegenerateResponse, error) {
	rawSecret := fmt.Sprintf("sk_live_%s", generateSecureRandom(32))
	hash := sha256.Sum256([]byte(rawSecret))
	hashedSecret := hex.EncodeToString(hash[:])

	if err := s.repo.UpdateSecret(ctx, id, userID, hashedSecret); err != nil {
		return nil, err
	}

	// Fetch updated key to get name and public key
	keys, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var key *models.ApiKey
	for _, k := range keys {
		if k.ID == id {
			key = &k
			break
		}
	}

	if key == nil {
		return nil, fmt.Errorf("API key not found")
	}

	return &models.ApiKeyRegenerateResponse{
		ID:        key.ID,
		Name:      key.Name,
		PublicKey: key.PublicKey,
		SecretKey: rawSecret,
		UpdatedAt: time.Now(),
	}, nil
}

func generateSecureRandom(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

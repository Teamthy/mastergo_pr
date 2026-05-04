package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
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
	log.Printf("CreateKey service: generating key for user %s with name %s", userID, name)

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

	log.Printf("CreateKey service: saving API key %s to database", apiKey.ID)
	if err := s.repo.Create(ctx, apiKey); err != nil {
		log.Printf("CreateKey service: ERROR saving to database: %v", err)
		return nil, err
	}

	log.Printf("CreateKey service: API key %s created successfully", apiKey.ID)
	return &models.ApiKeyCreateResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		PublicKey: apiKey.PublicKey,
		SecretKey: rawSecret,
		CreatedAt: apiKey.CreatedAt,
	}, nil
}

func (s *ApiKeyService) ListKeys(ctx context.Context, userID uuid.UUID) ([]models.ApiKey, error) {
	log.Printf("ListKeys service: fetching keys for user %s", userID)
	keys, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		log.Printf("ListKeys service: ERROR fetching keys: %v", err)
		return nil, err
	}
	log.Printf("ListKeys service: found %d keys for user %s", len(keys), userID)
	return keys, nil
}

func (s *ApiKeyService) RevokeKey(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	log.Printf("RevokeKey service: revoking key %s for user %s", id, userID)
	if err := s.repo.Revoke(ctx, id, userID); err != nil {
		log.Printf("RevokeKey service: ERROR revoking key: %v", err)
		return err
	}
	log.Printf("RevokeKey service: key %s revoked successfully", id)
	return nil
}

func (s *ApiKeyService) RegenerateKey(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.ApiKeyRegenerateResponse, error) {
	log.Printf("RegenerateKey service: regenerating key %s for user %s", id, userID)

	rawSecret := fmt.Sprintf("sk_live_%s", generateSecureRandom(32))
	hash := sha256.Sum256([]byte(rawSecret))
	hashedSecret := hex.EncodeToString(hash[:])

	if err := s.repo.UpdateSecret(ctx, id, userID, hashedSecret); err != nil {
		log.Printf("RegenerateKey service: ERROR updating secret: %v", err)
		return nil, err
	}

	// Fetch updated key to get name and public key
	keys, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		log.Printf("RegenerateKey service: ERROR fetching updated key: %v", err)
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
		log.Printf("RegenerateKey service: ERROR key not found after update")
		return nil, fmt.Errorf("API key not found")
	}

	log.Printf("RegenerateKey service: key %s regenerated successfully", id)
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

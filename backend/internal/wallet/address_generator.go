package wallet

import (
	"context"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

type AddressGenerator struct {
	enc  Encryptor
	repo WalletRepository
}

type Encryptor interface {
	Encrypt([]byte) (string, error)
}

type WalletRepository interface {
	SaveWallet(ctx context.Context, w *Wallet) error
}

type Wallet struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Address      string
	EncryptedKey string
}

func NewAddressGenerator(enc Encryptor, repo WalletRepository) *AddressGenerator {
	return &AddressGenerator{enc: enc, repo: repo}
}

func (g *AddressGenerator) CreateWallet(ctx context.Context, userID uuid.UUID) (*Wallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	privateBytes := crypto.FromECDSA(privateKey)

	encKey, err := g.enc.Encrypt(privateBytes)
	if err != nil {
		return nil, err
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	wallet := &Wallet{
		ID:           uuid.New(),
		UserID:       userID,
		Address:      address,
		EncryptedKey: encKey,
	}

	if err := g.repo.SaveWallet(ctx, wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

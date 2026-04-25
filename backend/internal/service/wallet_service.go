package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	"backend/internal/repository"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

var (
	ErrInsufficientBalance       = errors.New("insufficient balance")
	ErrInvalidDestinationAddress = errors.New("invalid destination address")
	ErrInvalidAmountWei          = errors.New("invalid amount wei")
	ErrWalletKeyMismatch         = errors.New("wallet key does not match stored address")
)

type WalletService struct {
	repo repository.WalletRepository
	key  []byte
	eth  *ethclient.Client
}

func NewWalletService(repo repository.WalletRepository, masterKey []byte, rpcURL string) (*WalletService, error) {
	if len(masterKey) != 32 {
		return nil, fmt.Errorf("invalid master key length: got %d, want 32", len(masterKey))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ethereum rpc: %w", err)
	}

	return &WalletService{
		repo: repo,
		key:  masterKey,
		eth:  client,
	}, nil
}

func (s *WalletService) Close() {
	if s.eth != nil {
		s.eth.Close()
	}
}

func (s *WalletService) EthClient() *ethclient.Client {
	return s.eth
}

func (s *WalletService) encrypt(data []byte) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(crand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return hex.EncodeToString(ciphertext), nil
}

func (s *WalletService) decrypt(enc string) ([]byte, error) {
	data, err := hex.DecodeString(enc)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("invalid ciphertext")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (s *WalletService) GenerateEthereumWallet() (string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	privateKeyBytes := crypto.FromECDSA(privateKey)

	encryptedKey, err := s.encrypt(privateKeyBytes)
	if err != nil {
		return "", "", err
	}

	return address, encryptedKey, nil
}

func (s *WalletService) CreateEthWallet(ctx context.Context, userID uuid.UUID) (string, error) {
	address, encryptedKey, err := s.GenerateEthereumWallet()
	if err != nil {
		return "", err
	}

	if err := s.repo.CreateWallet(ctx, userID, address, encryptedKey); err != nil {
		return "", err
	}

	return address, nil
}

func (s *WalletService) GetBalance(ctx context.Context, userID uuid.UUID) (string, error) {
	return s.repo.GetBalanceForUpdate(ctx, userID)
}

func (s *WalletService) GetHistory(ctx context.Context, userID uuid.UUID) ([]repository.TxHistory, error) {
	return s.repo.GetHistory(ctx, userID)
}

func parseWei(amountWei string) (*big.Int, error) {
	amountWei = strings.TrimSpace(amountWei)
	if amountWei == "" {
		return nil, ErrInvalidAmountWei
	}

	v, ok := new(big.Int).SetString(amountWei, 10)
	if !ok || v.Sign() <= 0 {
		return nil, ErrInvalidAmountWei
	}

	return v, nil
}

func (s *WalletService) Withdraw(
	ctx context.Context,
	userID uuid.UUID,
	to string,
	amountWei string,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if !common.IsHexAddress(to) {
		return "", ErrInvalidDestinationAddress
	}

	amount, err := parseWei(amountWei)
	if err != nil {
		return "", err
	}

	balanceStr, err := s.repo.GetBalanceForUpdate(ctx, userID)
	if err != nil {
		return "", err
	}

	balance := new(big.Int)
	if _, ok := balance.SetString(balanceStr, 10); !ok {
		return "", errors.New("invalid stored balance")
	}

	if balance.Cmp(amount) < 0 {
		return "", ErrInsufficientBalance
	}

	walletAddress, encKey, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return "", err
	}

	privBytes, err := s.decrypt(encKey)
	if err != nil {
		return "", err
	}

	privateKey, err := crypto.ToECDSA(privBytes)
	if err != nil {
		return "", err
	}

	derived := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	if !strings.EqualFold(derived, walletAddress) {
		return "", ErrWalletKeyMismatch
	}

	from := common.HexToAddress(walletAddress)
	toAddr := common.HexToAddress(to)

	nonce, err := s.eth.PendingNonceAt(ctx, from)
	if err != nil {
		return "", err
	}

	chainID, err := s.eth.NetworkID(ctx)
	if err != nil {
		return "", err
	}

	header, err := s.eth.HeaderByNumber(ctx, nil)
	if err != nil {
		return "", err
	}

	var tx *types.Transaction

	if header.BaseFee != nil {
		tip, err := s.eth.SuggestGasTipCap(ctx)
		if err != nil {
			return "", err
		}

		feeCap := new(big.Int).Add(
			new(big.Int).Mul(header.BaseFee, big.NewInt(2)),
			tip,
		)

		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			To:        &toAddr,
			Value:     amount,
			Gas:       21000,
			GasTipCap: tip,
			GasFeeCap: feeCap,
		})

		tx, err = types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
		if err != nil {
			return "", err
		}
	} else {
		gasPrice, err := s.eth.SuggestGasPrice(ctx)
		if err != nil {
			return "", err
		}

		tx = types.NewTransaction(
			nonce,
			toAddr,
			amount,
			21000,
			gasPrice,
			nil,
		)

		tx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			return "", err
		}
	}

	if err := s.eth.SendTransaction(ctx, tx); err != nil {
		return "", err
	}

	txHash := tx.Hash().Hex()

	if err := s.repo.MarkBroadcasted(ctx, userID, amountWei, to, txHash); err != nil {
		return "", err
	}

	return txHash, nil
}

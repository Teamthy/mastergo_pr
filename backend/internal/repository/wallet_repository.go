package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BroadcastedWithdrawal struct {
	TxHash string `json:"tx_hash"`
}

type TxHistory struct {
	TxHash    string    `json:"tx_hash"`
	Type      string    `json:"type"`
	AmountWei string    `json:"amount_wei"`
	To        string    `json:"to"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type WalletRepository interface {
	CreateWallet(ctx context.Context, userID uuid.UUID, address, encryptedKey string) error
	GetWalletByUserID(ctx context.Context, userID uuid.UUID) (address, encryptedKey string, err error)

	GetBalanceForUpdate(ctx context.Context, userID uuid.UUID) (string, error)
	MarkBroadcasted(ctx context.Context, userID uuid.UUID, amountWei, to, txHash string) error

	GetPendingBroadcasts(ctx context.Context) ([]BroadcastedWithdrawal, error)
	MarkConfirmed(ctx context.Context, txHash string) error
	MarkFailedAndRefund(ctx context.Context, txHash string) error

	GetHistory(ctx context.Context, userID uuid.UUID) ([]TxHistory, error)
}

type walletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) CreateWallet(
	ctx context.Context,
	userID uuid.UUID,
	address, encryptedKey string,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO wallets (user_id, address, encrypted_private_key, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO NOTHING
	`,
		userID,
		address,
		encryptedKey,
		time.Now(),
	)
	return err
}

func (r *walletRepository) GetWalletByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (string, string, error) {
	var address, encrypted string

	err := r.db.QueryRow(ctx, `
		SELECT address, encrypted_private_key
		FROM wallets
		WHERE user_id = $1
	`, userID).Scan(&address, &encrypted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", errors.New("wallet not found")
		}
		return "", "", err
	}

	return address, encrypted, nil
}

func (r *walletRepository) GetBalanceForUpdate(
	ctx context.Context,
	userID uuid.UUID,
) (string, error) {
	var balance string

	err := r.db.QueryRow(ctx, `
		SELECT balance::text
		FROM balances
		WHERE user_id = $1
	`, userID).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("balance not found")
		}
		return "", err
	}

	return balance, nil
}

func (r *walletRepository) MarkBroadcasted(
	ctx context.Context,
	userID uuid.UUID,
	amountWei, to, txHash string,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	res, err := tx.Exec(ctx, `
		INSERT INTO transactions (
			id, user_id, type, amount, to_address, tx_hash, status, created_at
		)
		VALUES ($1, $2, 'withdraw', $3::numeric, $4, $5, 'broadcasted', $6)
		ON CONFLICT (tx_hash) DO NOTHING
	`,
		uuid.New(),
		userID,
		amountWei,
		to,
		txHash,
		time.Now(),
	)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return tx.Commit(ctx)
	}

	balRes, err := tx.Exec(ctx, `
		UPDATE balances
		SET balance = balance - $1::numeric
		WHERE user_id = $2
		  AND balance >= $1::numeric
	`, amountWei, userID)
	if err != nil {
		return err
	}

	balRows := balRes.RowsAffected()
	if balRows == 0 {
		return errors.New("insufficient balance")
	}

	return tx.Commit(ctx)
}

func (r *walletRepository) GetPendingBroadcasts(
	ctx context.Context,
) ([]BroadcastedWithdrawal, error) {
	rows, err := r.db.Query(ctx, `
		SELECT tx_hash
		FROM transactions
		WHERE status = 'broadcasted'
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []BroadcastedWithdrawal

	for rows.Next() {
		var tx BroadcastedWithdrawal
		if err := rows.Scan(&tx.TxHash); err != nil {
			return nil, err
		}
		list = append(list, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *walletRepository) MarkConfirmed(
	ctx context.Context,
	txHash string,
) error {
	_, err := r.db.Exec(ctx, `
		UPDATE transactions
		SET status = 'confirmed'
		WHERE tx_hash = $1
	`, txHash)
	return err
}

func (r *walletRepository) MarkFailedAndRefund(
	ctx context.Context,
	txHash string,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var userID uuid.UUID
	var amount string

	err = tx.QueryRow(ctx, `
		SELECT user_id, amount::text
		FROM transactions
		WHERE tx_hash = $1
		FOR UPDATE
	`, txHash).Scan(&userID, &amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("transaction not found")
		}
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE balances
		SET balance = balance + $1::numeric
		WHERE user_id = $2
	`, amount, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE transactions
		SET status = 'refunded'
		WHERE tx_hash = $1
	`, txHash)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *walletRepository) GetHistory(
	ctx context.Context,
	userID uuid.UUID,
) ([]TxHistory, error) {
	rows, err := r.db.Query(ctx, `
		SELECT tx_hash, type, amount::text, to_address, status, created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TxHistory

	for rows.Next() {
		var tx TxHistory
		if err := rows.Scan(
			&tx.TxHash,
			&tx.Type,
			&tx.AmountWei,
			&tx.To,
			&tx.Status,
			&tx.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

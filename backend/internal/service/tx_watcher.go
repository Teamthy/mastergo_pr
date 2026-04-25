package service

import (
	"context"
	"errors"
	"log"
	"time"

	"backend/internal/repository"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TxWatcher struct {
	repo          repository.WalletRepository
	eth           *ethclient.Client
	interval      time.Duration
	confirmations uint64
}

func NewTxWatcher(repo repository.WalletRepository, eth *ethclient.Client) *TxWatcher {
	return &TxWatcher{
		repo:          repo,
		eth:           eth,
		interval:      15 * time.Second,
		confirmations: 6,
	}
}

func (w *TxWatcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.scan(ctx)
		}
	}
}

func (w *TxWatcher) scan(ctx context.Context) {
	txs, err := w.repo.GetPendingBroadcasts(ctx)
	if err != nil {
		log.Printf("tx watcher fetch error: %v", err)
		return
	}

	if len(txs) == 0 {
		return
	}

	head, err := w.eth.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Printf("tx watcher header error: %v", err)
		return
	}

	latestNumber := head.Number.Uint64()

	for _, tx := range txs {
		receipt, err := w.eth.TransactionReceipt(ctx, common.HexToHash(tx.TxHash))
		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				continue
			}
			log.Printf("tx watcher receipt error tx=%s err=%v", tx.TxHash, err)
			continue
		}

		if receipt == nil || receipt.BlockNumber == nil {
			continue
		}

		if latestNumber < receipt.BlockNumber.Uint64() {
			continue
		}

		confirmations := latestNumber - receipt.BlockNumber.Uint64()
		if confirmations < w.confirmations {
			continue
		}

		if receipt.Status == types.ReceiptStatusSuccessful {
			if err := w.repo.MarkConfirmed(ctx, tx.TxHash); err != nil {
				log.Printf("tx watcher confirm update error tx=%s err=%v", tx.TxHash, err)
			}
			continue
		}

		if err := w.repo.MarkFailedAndRefund(ctx, tx.TxHash); err != nil {
			log.Printf("tx watcher refund update error tx=%s err=%v", tx.TxHash, err)
		}
	}
}

package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/service"
)

type WalletHandler struct {
	service *service.WalletService
}

func NewWalletHandler(s *service.WalletService) *WalletHandler {
	return &WalletHandler{service: s}
}

type createWalletResponse struct {
	Address string `json:"address"`
	Network string `json:"network"`
}

type balanceResponse struct {
	BalanceWei string `json:"balance_wei"`
	Symbol     string `json:"symbol"`
}

type withdrawRequest struct {
	AmountWei string `json:"amount_wei"`
	To        string `json:"to"`
}

type withdrawResponse struct {
	TxHash string `json:"tx_hash"`
	Status string `json:"status"`
}

// POST /wallet/create
func (h *WalletHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	address, err := h.service.CreateEthWallet(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, createWalletResponse{
		Address: address,
		Network: "ethereum",
	})
}

// GET /wallet/balance
func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	balance, err := h.service.GetBalance(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, balanceResponse{
		BalanceWei: balance,
		Symbol:     "ETH",
	})
}

// GET /wallet/transactions
func (h *WalletHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	txs, err := h.service.GetHistory(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, txs)
}

// POST /wallet/withdraw
func (h *WalletHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req withdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if req.To == "" || req.AmountWei == "" {
		http.Error(w, "to and amount_wei are required", http.StatusBadRequest)
		return
	}

	txHash, err := h.service.Withdraw(r.Context(), userID, req.To, req.AmountWei)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, withdrawResponse{
		TxHash: txHash,
		Status: "broadcasted",
	})
}

package httphandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"test_project/model"
	"test_project/repository"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ServiceInt interface {
	GetBalance(ctx context.Context, walletUUID uuid.UUID) (balance int64, err error)
	Deposit(ctx context.Context, walletUUID uuid.UUID, amount int64) (int64, error)
	Withdraw(ctx context.Context, walletUUID uuid.UUID, amount int64) (int64, error)
}
type HTTPHandlers struct {
	WalletService ServiceInt
}

func NewHTTPHandlers(WalletService ServiceInt) *HTTPHandlers {
	return &HTTPHandlers{
		WalletService: WalletService,
	}
}

func (h *HTTPHandlers) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	var req UpdateWalletRequest
	var balance int64
	var err error
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "")
		return
	}
	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "Сумма не может быть меньше 0")
		return
	}
	if req.OperationType != string(model.Deposit) && req.OperationType != string(model.Withdraw) {
		writeError(w, http.StatusBadRequest, "Тип операции не верный")
		return
	}
	if req.OperationType == string(model.Deposit) {
		balance, err = h.WalletService.Deposit(r.Context(), req.WalletID, req.Amount)
	} else if req.OperationType == string(model.Withdraw) {
		balance, err = h.WalletService.Withdraw(r.Context(), req.WalletID, req.Amount)
	}
	if err != nil {
		if errors.Is(err, repository.ErrWalletNotFound) {
			writeError(w, http.StatusNotFound, "Не найден кошелек")
			return
		} else if errors.Is(err, repository.ErrInsufficientFunds) {
			writeError(w, http.StatusConflict, "Не достаточно средств")
			return
		} else {
			writeError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}
	}

	resp := UpdateWalletResponse{
		WalletID:      req.WalletID,
		OperationType: req.OperationType,
		Amount:        req.Amount,
		Balance:       balance,
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *HTTPHandlers) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletIDStr := vars["WALLET_UUID"]
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Не верный wallet id")
		return
	}
	balance, err := h.WalletService.GetBalance(r.Context(), walletID)
	if err != nil {
		if errors.Is(err, repository.ErrWalletNotFound) {
			writeError(w, http.StatusNotFound, "Не найден кошелёк")
			return
		} else {
			writeError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}
	}
	resp := GetBalanceResponse{
		WalletID: walletID,
		Balance:  balance,
	}
	writeJSON(w, http.StatusOK, resp)

}

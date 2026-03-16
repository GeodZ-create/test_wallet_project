package httphandlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	httphandlers "test_project/httpHandlers"
	"test_project/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type fakeService struct {
	getBalanceResult int64
	getBalanceErr    error

	depositResult int64
	depositErr    error

	withdrawResult int64
	withdrawErr    error
}

func (f *fakeService) GetBalance(ctx context.Context, walletUUID uuid.UUID) (int64, error) {
	return f.getBalanceResult, f.getBalanceErr
}

func (f *fakeService) Deposit(ctx context.Context, walletUUID uuid.UUID, amount int64) (int64, error) {
	return f.depositResult, f.depositErr
}

func (f *fakeService) Withdraw(ctx context.Context, walletUUID uuid.UUID, amount int64) (int64, error) {
	return f.withdrawResult, f.withdrawErr
}

func TestGetBalanceInvalidUUID(t *testing.T) {
	svc := &fakeService{}
	handler := httphandlers.NewHTTPHandlers(svc)

	req := httptest.NewRequest("GET", "/api/v1/wallets/not-a-uuid", nil)

	req = mux.SetURLVars(req, map[string]string{
		"WALLET_UUID": "not-a-uuid",
	})

	rec := httptest.NewRecorder()

	handler.GetBalance(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatal("Неверный статус")
	}
}

func TestGetBalanceNotFound(t *testing.T) {
	svc := &fakeService{
		getBalanceErr: repository.ErrWalletNotFound,
	}
	handler := httphandlers.NewHTTPHandlers(svc)
	walletID := uuid.New()

	req := httptest.NewRequest("GET", "/api/v1/wallets/uuid", nil)
	req = mux.SetURLVars(req, map[string]string{
		"WALLET_UUID": walletID.String(),
	})
	rec := httptest.NewRecorder()

	handler.GetBalance(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatal("Не верный статус")
	}
}

func TestDepositSuccess(t *testing.T) {
	walletID := uuid.New()
	svc := &fakeService{
		depositResult: 150,
	}
	handler := httphandlers.NewHTTPHandlers(svc)
	fakeBody := &httphandlers.UpdateWalletRequest{
		WalletID:      walletID,
		OperationType: "deposit",
		Amount:        100,
	}
	fakeBodyJson, err := json.Marshal(fakeBody)
	if err != nil {
		t.Fatal("Ошибка парсинга JSON")
	}

	req := httptest.NewRequest("POST", "/api/v1/wallet", bytes.NewReader(fakeBodyJson))
	rec := httptest.NewRecorder()

	handler.UpdateWallet(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatal("Ошибка статуса", rec.Code, http.StatusOK)
	}

}

func TestWithdrawInsufficientFunds(t *testing.T) {
	svc := &fakeService{
		withdrawErr: repository.ErrInsufficientFunds,
	}
	walletID := uuid.New()
	handler := httphandlers.NewHTTPHandlers(svc)
	fakeBody := &httphandlers.UpdateWalletRequest{
		WalletID:      walletID,
		OperationType: "withdraw",
		Amount:        100,
	}
	fakeBodyJson, err := json.Marshal(fakeBody)
	if err != nil {
		t.Fatal("Ошибка парсинга JSON")
	}

	req := httptest.NewRequest("POST", "/api/v1/wallet", bytes.NewReader(fakeBodyJson))
	rec := httptest.NewRecorder()

	handler.UpdateWallet(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatal("Ошибка статуса")
	}
}

func TestInvalidJson(t *testing.T) {
	svc := &fakeService{}
	handler := httphandlers.NewHTTPHandlers(svc)
	fakeJsonString := "sadasdsa"
	req := httptest.NewRequest("POST", "/api/v1/wallet", bytes.NewBufferString(fakeJsonString))
	rec := httptest.NewRecorder()

	handler.UpdateWallet(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatal("Неверный статус", rec.Code, http.StatusBadRequest)
	}
}

func TestInternalError(t *testing.T) {
	svc := &fakeService{
		depositErr: errors.New("Ошибка"),
	}
	handler := httphandlers.NewHTTPHandlers(svc)
	walletID := uuid.New()
	fakeBody := &httphandlers.UpdateWalletRequest{
		WalletID:      walletID,
		OperationType: "deposit",
		Amount:        100,
	}
	fakeBodyJson, err := json.Marshal(fakeBody)
	if err != nil {
		t.Fatal("Ошибка парсинга JSON")
	}

	req := httptest.NewRequest("POST", "/api/v1/wallet", bytes.NewReader(fakeBodyJson))
	rec := httptest.NewRecorder()

	handler.UpdateWallet(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatal("Неверный статус ответа", rec.Code, http.StatusInternalServerError)
	}

}

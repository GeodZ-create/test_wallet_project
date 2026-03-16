package httphandlers

import "github.com/google/uuid"

type UpdateWalletRequest struct {
	WalletID      uuid.UUID `json:"walletId"`
	OperationType string    `json:"operationType"`
	Amount        int64     `json:"amount"`
}

type GetBalanceResponse struct {
	WalletID uuid.UUID `json:"walletId"`
	Balance  int64     `json:"balance"`
}

type UpdateWalletResponse struct {
	WalletID      uuid.UUID `json:"walletId"`
	OperationType string    `json:"operationType"`
	Amount        int64     `json:"amount"`
	Balance       int64     `json:"balance"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

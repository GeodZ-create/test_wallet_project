package service

import (
	"context"
	"test_project/repository"

	"github.com/google/uuid"
)

type WalletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) *WalletService {
	return &WalletService{
		repo: repo,
	}
}

func (s *WalletService) GetBalance(ctx context.Context, walletUUID uuid.UUID) (balance int64, err error) {
	balance, err = s.repo.GetBalance(ctx, walletUUID)
	return balance, err
}

func (s *WalletService) Deposit(ctx context.Context, walletUUID uuid.UUID, amount int64) (int64, error) {
	balance, err := s.repo.Deposit(ctx, walletUUID, amount)
	return balance, err
}
func (s *WalletService) Withdraw(ctx context.Context, walletUUID uuid.UUID, amount int64) (int64, error) {
	balance, err := s.repo.Withdraw(ctx, walletUUID, amount)
	return balance, err
}

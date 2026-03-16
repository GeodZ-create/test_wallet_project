package repository

import (
	"context"
	"errors"
	"fmt"
	"test_project/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	GetBalance(ctx context.Context, walletID uuid.UUID) (balance int64, err error)
	Deposit(ctx context.Context, walletID uuid.UUID, amount int64) (int64, error)
	Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) (int64, error)
}

type DbRepo struct {
	db *gorm.DB
}

func NewDbRepo(db *gorm.DB) *DbRepo {
	return &DbRepo{
		db: db,
	}
}

var ErrWalletNotFound = errors.New("Запись не найдена")
var ErrInsufficientFunds = errors.New("Недостаточно средств")

func (r *DbRepo) GetBalance(ctx context.Context, walletID uuid.UUID) (balance int64, err error) {
	var wallet model.Wallet
	err = r.db.WithContext(ctx).Where("wallet_id = ?", walletID).First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, ErrWalletNotFound
	} else if err != nil {
		return 0, fmt.Errorf("Ошибка доступа к БД: %w", err)
	}
	return wallet.Balance, nil
}

func (r *DbRepo) Deposit(ctx context.Context, walletID uuid.UUID, amount int64) (int64, error) {
	var wallet model.Wallet
	result := r.db.WithContext(ctx).Model(&wallet).Clauses(clause.Returning{Columns: []clause.Column{{Name: "balance"}}}).Where("wallet_id = ?", walletID).Update("balance", gorm.Expr("balance + ?", amount))
	if result.Error != nil {
		return 0, fmt.Errorf("Ошибка доступа к БД: %w", result.Error)
	} else if result.RowsAffected == 0 {
		return 0, ErrWalletNotFound
	}
	return wallet.Balance, nil
}

func (r *DbRepo) Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) (int64, error) {
	var wallet model.Wallet
	result := r.db.WithContext(ctx).Model(&wallet).Clauses(clause.Returning{Columns: []clause.Column{{Name: "balance"}}}).Where("wallet_id = ? AND balance >= ?", walletID, amount).Update("balance", gorm.Expr("balance - ?", amount))
	if result.Error != nil {
		return 0, fmt.Errorf("Ошибка доступа к БД: %w", result.Error)
	} else if result.RowsAffected == 0 {
		var wallet model.Wallet
		err := r.db.WithContext(ctx).Where("wallet_id = ?", walletID).First(&wallet).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrWalletNotFound
		} else if err != nil {
			return 0, fmt.Errorf("Ошибка доступа к БД: %w", err)
		}
		return 0, ErrInsufficientFunds
	}
	return wallet.Balance, nil
}

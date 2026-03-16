package repository_test

import (
	"context"
	"errors"
	"os"
	"sync"
	"test_project/model"
	"test_project/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDbRepoDeposit(t *testing.T) {
	godotenv.Load("../config.env")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL is required")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("Проблема с доступом к БД", err)
	}
	repo := repository.NewDbRepo(db)
	err = db.Exec("TRUNCATE TABLE wallets RESTART IDENTITY").Error
	if err != nil {
		t.Fatal("Не получилось очистить таблицу", err)
	}
	walletID := uuid.New()
	testWallet := model.Wallet{
		WalletID: walletID,
		Balance:  0,
	}
	err = db.Create(&testWallet).Error
	if err != nil {
		t.Fatal("Не получилось создать кошелек", err)
	}
	balance, err := repo.Deposit(context.Background(), testWallet.WalletID, 100)
	if err != nil {
		t.Fatal("Не удалось получить выполнить депозит", err)
	}

	dbBalance, err := repo.GetBalance(context.Background(), testWallet.WalletID)
	if err != nil {
		t.Fatal("Не удалось получить баланс", err)
	}
	if balance != 100 || dbBalance != 100 {
		t.Fatal("Баланс не правильный", balance, dbBalance)
	}
}

func TestDbRepoWithdraw(t *testing.T) {
	godotenv.Load("../config.env")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL is required")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("Проблема с доступом к БД", err)
	}
	repo := repository.NewDbRepo(db)
	err = db.Exec("TRUNCATE TABLE wallets RESTART IDENTITY").Error
	if err != nil {
		t.Fatal("Не получилось очистить таблицу", err)
	}
	walletID := uuid.New()
	testWallet := model.Wallet{
		WalletID: walletID,
		Balance:  100,
	}
	err = db.Create(&testWallet).Error
	if err != nil {
		t.Fatal("Не получилось создать кошелек", err)
	}
	balance, err := repo.Withdraw(context.Background(), testWallet.WalletID, 40)
	if err != nil {
		t.Fatal("Не удалось получить выполнить депозит", err)
	}

	dbBalance, err := repo.GetBalance(context.Background(), testWallet.WalletID)
	if err != nil {
		t.Fatal("Не удалось получить баланс", err)
	}
	if balance != 60 || dbBalance != 60 {
		t.Fatal("Баланс не правильный", balance, dbBalance)
	}
}

func TestDbRepoWithdrawInsufficient(t *testing.T) {
	godotenv.Load("../config.env")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL is required")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("Проблема с доступом к БД", err)
	}
	repo := repository.NewDbRepo(db)
	err = db.Exec("TRUNCATE TABLE wallets RESTART IDENTITY").Error
	if err != nil {
		t.Fatal("Не получилось очистить таблицу", err)
	}
	walletID := uuid.New()
	testWallet := model.Wallet{
		WalletID: walletID,
		Balance:  100,
	}
	err = db.Create(&testWallet).Error
	if err != nil {
		t.Fatal("Не получилось создать кошелек", err)
	}
	_, err = repo.Withdraw(context.Background(), testWallet.WalletID, 200)
	if !errors.Is(err, repository.ErrInsufficientFunds) {
		t.Fatal("Не верная ошибка", err)
	}
	dbBalance, err := repo.GetBalance(context.Background(), testWallet.WalletID)
	if err != nil {
		t.Fatal("Не удалось получить баланс", err)
	}
	if dbBalance != 100 {
		t.Fatal("Баланс не правильный", dbBalance)
	}
}

func TestDbRepoGetBalance(t *testing.T) {
	godotenv.Load("../config.env")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL is required")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("Проблема с доступом к БД", err)
	}
	repo := repository.NewDbRepo(db)
	err = db.Exec("TRUNCATE TABLE wallets RESTART IDENTITY").Error
	if err != nil {
		t.Fatal("Не получилось очистить таблицу", err)
	}
	walletID := uuid.New()
	testWallet := model.Wallet{
		WalletID: walletID,
		Balance:  100,
	}
	err = db.Create(&testWallet).Error
	if err != nil {
		t.Fatal("Не получилось создать кошелек", err)
	}
	balance, err := repo.GetBalance(context.Background(), testWallet.WalletID)
	if err != nil {
		t.Fatal("Не удалось получить выполнить депозит", err)
	}
	if balance != 100 {
		t.Fatal("Баланс не правильный", balance)
	}
}

func TestDbRepoGetBalanceNotFound(t *testing.T) {
	godotenv.Load("../config.env")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL is required")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("Проблема с доступом к БД", err)
	}
	repo := repository.NewDbRepo(db)
	err = db.Exec("TRUNCATE TABLE wallets RESTART IDENTITY").Error
	if err != nil {
		t.Fatal("Не получилось очистить таблицу", err)
	}
	walletID := uuid.New()

	_, err = repo.GetBalance(context.Background(), walletID)
	if !errors.Is(err, repository.ErrWalletNotFound) {
		t.Fatal("Не верная ошибка", err)
	}
}

func TestDbRepoDepositNotFound(t *testing.T) {
	godotenv.Load("../config.env")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL is required")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("Проблема с доступом к БД", err)
	}
	repo := repository.NewDbRepo(db)
	err = db.Exec("TRUNCATE TABLE wallets RESTART IDENTITY").Error
	if err != nil {
		t.Fatal("Не получилось очистить таблицу", err)
	}
	walletId := uuid.New()
	_, err = repo.Deposit(context.Background(), walletId, 40)
	if !errors.Is(err, repository.ErrWalletNotFound) {
		t.Fatal("Не верная ошибка", err)
	}

}

func TestDbRepoConcurrentDeposits(t *testing.T) {
	godotenv.Load("../config.env")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL is required")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("Проблема с доступом к БД", err)
	}
	repo := repository.NewDbRepo(db)
	err = db.Exec("TRUNCATE TABLE wallets RESTART IDENTITY").Error
	if err != nil {
		t.Fatal("Не получилось очистить таблицу", err)
	}
	walletID := uuid.New()
	testWallet := model.Wallet{
		WalletID: walletID,
		Balance:  0,
	}
	err = db.Create(&testWallet).Error
	if err != nil {
		t.Fatal("Не получилось создать кошелек", err)
	}
	var wg sync.WaitGroup
	errCh := make(chan error, 1000)
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := repo.Deposit(context.TODO(), testWallet.WalletID, 1)
			if err != nil {
				errCh <- err
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatal("Ошибка депозита", err)
		}
	}
	balance, err := repo.GetBalance(context.Background(), testWallet.WalletID)
	if err != nil {
		t.Fatal("Не удалось получить баланс", err)
	}

	if balance != 200 {
		t.Fatal("Случилась гонка данных")
	}

}

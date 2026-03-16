package model

import (
	"time"

	"github.com/google/uuid"
)

type OperationType string

const (
	Deposit  OperationType = "deposit"
	Withdraw OperationType = "withdraw"
)

type Wallet struct {
	ID        uint      `gorm:"primaryKey"`
	WalletID  uuid.UUID `gorm:"type:uuid; uniqueIndex;not null"`
	Balance   int64     `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

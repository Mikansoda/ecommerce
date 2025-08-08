package entity

import (
	"time"
	"github.com/google/uuid"
)

// Entity table for payments
type Payment struct {
	ID              uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrderID         uuid.UUID  `gorm:"type:char(36);not null" json:"order_id"`
	PaymentType     string     `gorm:"type:enum('transfer','internal_balance');not null" json:"payment_type"`
	TransferNumber  string     `gorm:"type:varchar(255)" json:"transfer_number"`
	Status          string     `gorm:"type:enum('pending','verified','rejected');default:'pending'" json:"status"`
	VerifiedBy      *uuid.UUID `gorm:"type:char(36)" json:"verified_by,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`

	// Relations
	Order *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Admin *Admin `gorm:"foreignKey:VerifiedBy" json:"verified_by_admin,omitempty"`
}

// Entity table for topups
type Topup struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	UserID         uuid.UUID  `gorm:"type:char(36);not null" json:"user_id"`
	TransferNumber string     `gorm:"type:varchar(100)" json:"transfer_number"`
	Amount         float64    `gorm:"type:decimal(12,2);not null" json:"amount"`
	Status         string     `gorm:"type:enum('pending','verified','rejected');default:'pending'" json:"status"`
	VerifiedBy     *uuid.UUID `gorm:"type:char(36)" json:"verified_by,omitempty"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	VerifiedAt     *time.Time `json:"verified_at,omitempty"`

	// Relations
	User  *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Admin *Admin `gorm:"foreignKey:VerifiedBy" json:"verified_by_admin,omitempty"`
}
package entity

import (
	"time"
	"github.com/google/uuid"

)

// Entity table for payments
type Payment struct {
	ID              uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrderID         uuid.UUID  `gorm:"type:char(36);not null" json:"order_id"`
	TransferNumber  string     `gorm:"type:varchar(255)" json:"transfer_number"`
	PaymentType     string     `gorm:"type:enum('transfer','internal_balance');not null" json:"payment_type"`
	Status          string     `gorm:"type:enum('pending','verified','failed');default:'pending'" json:"status"`
	VerifiedBy     *uuid.UUID  `gorm:"type:char(36)" json:"verified_by,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	VerifiedAt     *time.Time  `json:"verified_at,omitempty"`

	// Relations
	Order          *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Admin          *UsersAndAdmins `gorm:"foreignKey:VerifiedBy" json:"verified_by_admin,omitempty"`
}
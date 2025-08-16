package entity

import (
	"time"
	"github.com/google/uuid"

)

// Entity table for orders
type Order struct {
	ID           uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	UserID       uuid.UUID  `gorm:"type:char(36);not null" json:"user_id"`
	AddressID    uuid.UUID  `gorm:"type:char(36);not null" json:"address_id"`
	Status       string     `gorm:"type:enum('pending','paid','shipped','completed','cancelled');default:'pending'" json:"status"`
	TotalAmount  float64    `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	OrderDate    time.Time  `gorm:"autoCreateTime" json:"order_date"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	ExpiredAt   *time.Time  `json:"expired_at"` // expired in one hour in case not paid
    PaymentID   *uuid.UUID  `gorm:"type:char(36)" json:"payment_id,omitempty"`

	// Relations
	User        *UsersAndAdmins       `gorm:"foreignKey:UserID" json:"-"`
	Address     *Address    `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	Payment     *Payment   `gorm:"foreignKey:PaymentID;references:ID" json:"payment,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;references:ID" json:"order_items,omitempty"`
}

// Entity table for items within an order
type OrderItem struct {
    ID           uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrderID      uuid.UUID  `gorm:"type:char(36);not null" json:"order_id"`
	ProductID    uint       `gorm:"not null" json:"product_id"`
	Quantity     int        `gorm:"not null" json:"quantity"`
	PriceAtOrder float64    `gorm:"type:decimal(12,2);not null" json:"price"`

	// Relations
	Order       *Order   `gorm:"foreignKey:OrderID" json:"-"`
	Product     *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
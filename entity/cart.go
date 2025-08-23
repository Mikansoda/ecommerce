package entity

import (
	"time"
	"github.com/google/uuid"
)

// Entity table for cart per user
type Cart struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	UserID    uuid.UUID  `gorm:"type:char(36);not null" json:"user_id"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	User     *Users      `gorm:"foreignKey:UserID" json:"-"`
	Items   []CartItem   `gorm:"foreignKey:CartID" json:"items,omitempty"`
}

// Entity table for items within each carts
type CartItem struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	CartID    uuid.UUID  `gorm:"type:char(36);not null" json:"cart_id"`
	ProductID uint       `gorm:"not null" json:"product_id"`
	Quantity  int        `gorm:"not null" json:"quantity"`

	// Relations
	Cart     *Cart    `gorm:"foreignKey:CartID" json:"-"`
	Product  *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}


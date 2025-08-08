package entity

import (
	"time"
	// "golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

// Entity table for users
type User struct {
    ID           uuid.UUID     `gorm:"type:char(36);primaryKey" json:"id"`
    Username     string        `gorm:"type:varchar(100);not null" json:"username"`
    Email        string        `gorm:"type:varchar(100);unique;not null" json:"email"`
    PasswordHash string        `gorm:"type:text;not null" json:"-"`
    CreatedAt    time.Time     `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time     `gorm:"autoUpdateTime" json:"updated_at"`

    // Relations
    OTPCode    []UserOTP       `gorm:"foreignKey:UserID" json:"otps,omitempty"`
    Addresses  []Address       `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
    Favorites  []UserFavorite  `gorm:"foreignKey:UserID" json:"favorites,omitempty"`
}

// Entity table for user OTP codes
type UserOTP struct {
    ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
    UserID    uuid.UUID `gorm:"type:char(36);not null" json:"user_id"`
    OTPCode   string    `gorm:"type:varchar(100);not null" json:"otp_code"`
    ExpiredAt time.Time `json:"expired_at"`
    IsUsed    bool      `gorm:"default:false" json:"is_used"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

    // Relation
    User User `gorm:"foreignKey:UserID" json:"-"` // The user the OTP belongs to, cannot be null, always has an owner
}

// Entity table for user adresses
type Address struct {
    ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
    UserID       uuid.UUID `gorm:"type:char(36);not null" json:"user_id"`
    ReceiverName string    `gorm:"type:varchar(100);not null" json:"receiver_name"`
    PhoneNumber  string    `gorm:"type:varchar(100);not null" json:"phone_number"`
    AddressLine  string    `gorm:"type:varchar(500);not null" json:"address_line"`
    PostalCode   string    `gorm:"type:varchar(100)" json:"postal_code"`
    IsDefault    bool      `gorm:"default:false" json:"is_default"`

    // Relation
    User User `gorm:"foreignKey:UserID" json:"-"` // The user the address belongs to, cannot be null, always has an owner
}

// Entity table for user balance
type UserBalance struct {
    ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
    Balance   float64   `gorm:"type:decimal(12,2);default:0" json:"balance"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

    // Relation
    User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ID;references:ID" json:"-"` // The user the balance belongs to, cannot be null, always has an owner
}

// Entity table for user favorite products
type UserFavorite struct {
    ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID    uuid.UUID `gorm:"type:char(36);not null" json:"user_id"`
    ProductID int       `gorm:"not null" json:"product_id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

    // Relations
    User    User    `gorm:"foreignKey:UserID" json:"-"`
    Product Product `gorm:"foreignKey:ProductID" json:"-"` // The user the favorited product belongs to, cannot be null, always has an owner
}

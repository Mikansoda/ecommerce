package entity

import (
	"time"
	// "golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

// Entity table for admins
type Admin struct {
    ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
    Username     string    `gorm:"type:varchar(100);not null" json:"username"`
	Email        string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:text;not null" json:"-"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Payments []Payment  `gorm:"foreignKey:VerifiedBy" json:"payments_verified,omitempty"`
	Topups   []Topup    `gorm:"foreignKey:VerifiedBy" json:"topups_verified,omitempty"`
	OTPCode  []AdminOTP `gorm:"foreignKey:AdminID" json:"otps,omitempty"`
}

// Entity table for user OTP codes
type AdminOTP struct {
   ID         uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	AdminID   uuid.UUID `gorm:"type:char(36);not null" json:"admin_id"`
	OTPCode   string    `gorm:"type:varchar(100);not null" json:"otp_code"`
	ExpiredAt time.Time `json:"expired_at"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	Admin *Admin `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
}
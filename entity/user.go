package entity

import (
	"time"
	// "golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
    "gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// Entity table for users
type UsersAndAdmins struct {
    ID           uuid.UUID     `gorm:"type:char(36);primaryKey" json:"id"`
    Username     string        `gorm:"type:varchar(100);unique;not null" json:"username"`
    Email        string        `gorm:"type:varchar(100);unique;not null" json:"email"`
    PasswordHash string        `gorm:"type:text;not null" json:"-"`
    Role         Role          `gorm:"type:enum('user','admin');default:'user';not null" json:"role"`
    IsActive     bool          `gorm:"default:false"`
    OTPHash           string   `gorm:"size:191"` // hash dari OTP terakhir
	OTPExpiresAt     *time.Time
    RefreshTokenHash  string   `gorm:"size:191"`
	RefreshExpiresAt *time.Time
    CreatedAt    time.Time     `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time     `gorm:"autoUpdateTime" json:"updated_at"`

    // Relations
    Addresses  []Address       `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
}

// Entity table for user adresses
type Address struct {
    ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
    UserID       uuid.UUID `gorm:"type:char(36);not null" json:"-"`
    ReceiverName string    `gorm:"type:varchar(100);not null" json:"receiver_name"`
    PhoneNumber  string    `gorm:"type:varchar(100);not null" json:"phone_number"`
    AddressLine  string    `gorm:"type:varchar(500);not null" json:"address_line"`
    City         string    `gorm:"type:varchar(100);not null" json:"city"`
    Province     string    `gorm:"type:varchar(100);not null" json:"province"`
    PostalCode   string    `gorm:"type:varchar(100)" json:"postal_code"` // using VARCHAR instead of INT to prevent deletion of 0 at start of postal code
    CreatedAt    time.Time     `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`  // for soft deletion

    // Relation
    User         UsersAndAdmins `gorm:"foreignKey:UserID" json:"-"` // The user the address belongs to, cannot be null, always has an owner
}

func (u *UsersAndAdmins) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
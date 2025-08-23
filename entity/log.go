package entity

import (
	"time"
	"github.com/google/uuid"
)

// Entity table for logging actions performed by users or admins
type ActionLog struct {
    ID         uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
    ActorType  string     `gorm:"type:enum('user','admin','unknown');not null" json:"actor_type"`
    ActorID   *uuid.UUID  `gorm:"type:char(36)" json:"actor_id"` // nullable (for logins/register)
    Action     string     `gorm:"type:varchar(255);not null" json:"action"`
    EntityType string     `gorm:"type:enum('users','addresses','auth','carts','categories','orders','payments','products','product_images');not null" json:"entity_type"`
    EntityID   string     `gorm:"type:varchar(100);not null" json:"entity_id"`
    CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`

    User      *Users      `gorm:"foreignKey:ActorID;references:ID" json:"-"`
}

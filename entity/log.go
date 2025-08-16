package entity

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"

)

// Entity table for logging actions performed by users or admins
type ActionLog struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	ActorType   string    `gorm:"type:enum('user','admin');not null" json:"actor_type"`
	ActorID     uuid.UUID `gorm:"type:char(36);not null" json:"actor_id"`
	Action      string    `gorm:"type:varchar(255);not null" json:"action"`
	EntityType  string    `gorm:"type:enum('users','products','orders','payments','topups','addresses','categories','carts','cart_items','order_items');not null" json:"entity_type"`
	EntityID    string    `gorm:"type:varchar(100);not null" json:"entity_id"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`  // for soft deletion

	// Relation
    User       *UsersAndAdmins `gorm:"foreignKey:ActorID" json:"-"`
}
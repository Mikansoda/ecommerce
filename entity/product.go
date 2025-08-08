package entity

import (
	"time"
	"github.com/google/uuid"
)

// Entity table for products
type Product struct {
	ID          uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string            `gorm:"type:varchar(255);not null" json:"name"`
	Description string            `gorm:"type:varchar(1000)" json:"description,omitempty"`
	Price       float64           `gorm:"type:decimal(12,2);not null" json:"price"`
	Stock       int               `gorm:"not null" json:"stock"`
	CreatedAt   time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Images     []ProductImage     `gorm:"foreignKey:ProductID" json:"images,omitempty"`          // one-to-many -> has many images, belongs to 1 product
	Categories []ProductCategory  `gorm:"many2many:product_categories_map" json:"categories"`    // many-to-many, many products can belong to many categories, joint via PIVOT table
	Ratings    []Rating           `gorm:"foreignKey:ProductID" json:"ratings,omitempty"`         // one-to-many -> has many ratings, belongs to 1 user
}

// Entity table for producst images
type ProductImage struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	ImageURL  string    `gorm:"type:text;not null" json:"image_url"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`

	// Relations
	Product   *Product  `gorm:"foreignKey:ProductID" json:"-"` 
}

// Entity table for producst categories
type ProductCategory struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Products []Product `gorm:"many2many:product_categories_map" json:"products,omitempty"`
}

// Entity table for producst ratings
type Rating struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uuid.UUID      `gorm:"type:char(36);not null" json:"user_id"`
	ProductID uint           `gorm:"not null" json:"product_id"`
	Rating    int            `gorm:"type:int;check:rating >= 1 AND rating <= 5" json:"rating"`
	Review    string         `gorm:"type:varchar(100)" json:"review,omitempty"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	User    *User    `gorm:"foreignKey:UserID" json:"-"`
	Product *Product `gorm:"foreignKey:ProductID" json:"-"`
}
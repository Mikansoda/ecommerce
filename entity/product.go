package entity

import (
	"time"
	"gorm.io/gorm"

)

// Entity table for products
type Product struct {
	ID            uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string          `gorm:"type:varchar(255);not null" json:"name"`
	Description   string          `gorm:"type:varchar(1000)" json:"description,omitempty"`
	Price         float64         `gorm:"type:decimal(12,2);not null" json:"price"`
	Stock         uint            `gorm:"not null" json:"stock"`
	ExpiryYear    *int            `gorm:"type:int" json:"expiry_year,omitempty"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"autoUpdateTime" json:"updated_at"` 
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`  // for soft deletion

	// Relations
	Images      []ProductImage     `gorm:"foreignKey:ProductID" json:"images,omitempty"`          // one-to-many -> has many images, belongs to 1 product
	Categories  []ProductCategory  `gorm:"many2many:product_categories_map" json:"categories"`    // many-to-many, many products can belong to many categories, joint via PIVOT table
}

// Entity table for producst images
type ProductImage struct {
	ID           uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID    uint             `gorm:"not null" json:"product_id"`
	ImageURL     string           `gorm:"type:text;not null" json:"image_url"`
	IsPrimary    bool             `gorm:"default:false" json:"is_primary"`
	CreatedAt    time.Time        `gorm:"autoCreateTime" json:"created_at"`
    DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`  // for soft deletion

	// Relations
	Product     *Product          `gorm:"foreignKey:ProductID" json:"-"` 
}

// Entity table for producst categories
type ProductCategory struct {
	ID           uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string           `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt    time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`  // for soft deletion

	// Relations
	Products   []Product          `gorm:"many2many:product_categories_map" json:"products,omitempty"`
}

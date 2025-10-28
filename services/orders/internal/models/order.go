package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPaid      OrderStatus = "paid"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP;constraint:OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index;foreignKey:UserID;references:users(id);constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user_id"`
	TotalAmount float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status      OrderStatus    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Items       pq.StringArray `gorm:"type:text[]" json:"items"` // JSON строки: ["item1", "item2"]
}

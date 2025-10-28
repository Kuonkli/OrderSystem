package repository

import (
	"OrderSystem/pkg/logger"
	"OrderSystem/services/orders/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) (*models.Order, error)
	FindByID(id string) (*models.Order, error)
	FindByUserID(userID string, limit, offset int, sortBy, sortOrder string) ([]*models.Order, error)
	Update(order *models.Order) error
}

type orderRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewOrderRepository(db *gorm.DB, log *logger.Logger) OrderRepository {
	return &orderRepository{db: db, log: log}
}

func (r *orderRepository) Create(order *models.Order) (*models.Order, error) {
	result := r.db.Create(order)
	if result.Error != nil {
		r.log.Errorf("Failed to create order: %v", result.Error)
		return nil, result.Error
	}
	return order, nil
}

func (r *orderRepository) FindByID(id string) (*models.Order, error) {
	var order models.Order
	result := r.db.First(&order, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(userID string, limit, offset int, sortBy, sortOrder string) ([]*models.Order, error) {
	var orders []*models.Order

	// Допустимые поля
	validSortFields := map[string]string{
		"total_amount": "orders.total_amount",
		"created_at":   "orders.created_at",
		"status":       "orders.status",
	}

	// По умолчанию
	orderClause := "orders.created_at DESC"
	if field, ok := validSortFields[sortBy]; ok {
		if sortOrder == "ASC" {
			orderClause = field + " ASC"
		} else {
			orderClause = field + " DESC"
		}
	}

	result := r.db.Table("orders").
		Select("orders.*, users.email as user_email").
		Joins("JOIN users ON users.id = orders.user_id").
		Where("orders.user_id = ?", userID).
		Order(orderClause).
		Offset(offset).
		Limit(limit).
		Scan(&orders)
	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

func (r *orderRepository) Update(order *models.Order) error {
	result := r.db.Save(order)
	if result.Error != nil {
		r.log.Errorf("Failed to update order: %v", result.Error)
		return result.Error
	}
	return nil
}

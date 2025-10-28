package service

import (
	"OrderSystem/pkg/dto"
	"OrderSystem/pkg/logger"
	"OrderSystem/services/orders/internal/models"
	"OrderSystem/services/orders/internal/repository"
	"github.com/google/uuid"
)

type OrdersService struct {
	orderRepo repository.OrderRepository
	log       *logger.Logger
}

func NewOrdersService(orderRepo repository.OrderRepository, log *logger.Logger) *OrdersService {
	return &OrdersService{
		orderRepo: orderRepo,
		log:       log,
	}
}

func (s *OrdersService) CreateOrder(userID string, payload dto.CreateOrderRequest) (*models.Order, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	order := &models.Order{
		UserID:      userUUID,
		TotalAmount: payload.TotalAmount,
		Status:      models.StatusPending,
		Items:       payload.Items,
	}

	return s.orderRepo.Create(order)
}

func (s *OrdersService) GetOrder(orderID string) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	return &dto.OrderResponse{
		ID:          order.ID.String(),
		UserID:      order.UserID.String(),
		TotalAmount: order.TotalAmount,
		Status:      string(order.Status),
		Items:       order.Items,
		CreatedAt:   order.CreatedAt,
	}, nil
}

func (s *OrdersService) ListUserOrders(userID string, page, limit int, sortBy, sortOrder string) (*dto.OrderListResponse, error) {
	offset := page * limit
	orders, err := s.orderRepo.FindByUserID(userID, limit, offset, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	var items []dto.OrderResponse
	for _, o := range orders {
		items = append(items, dto.OrderResponse{
			ID:          o.ID.String(),
			UserID:      o.UserID.String(),
			TotalAmount: o.TotalAmount,
			Status:      string(o.Status),
			Items:       o.Items,
			CreatedAt:   o.CreatedAt,
		})
	}

	return &dto.OrderListResponse{Orders: items}, nil
}

func (s *OrdersService) UpdateOrderStatus(orderID string, status models.OrderStatus) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	order.Status = status
	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return &dto.OrderResponse{
		ID:          order.ID.String(),
		UserID:      order.UserID.String(),
		TotalAmount: order.TotalAmount,
		Status:      string(order.Status),
		Items:       order.Items,
		CreatedAt:   order.CreatedAt,
	}, nil
}

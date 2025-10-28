package dto

import (
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfile struct {
	ID        string   `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateOrderRequest struct {
	TotalAmount float64  `json:"total_amount" binding:"required"`
	Items       []string `json:"items" binding:"required"`
}

type OrderResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	Items       []string  `json:"items"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
}

package handlers

import (
	"OrderSystem/pkg/dto"
	"OrderSystem/pkg/logger"
	"OrderSystem/services/orders/internal/models"
	"OrderSystem/services/orders/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	orderService *service.OrdersService
	logger       *logger.Logger
}

func NewOrderHandler(orderService *service.OrdersService, logger *logger.Logger) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		logger:       logger,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	order, err := h.orderService.CreateOrder(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    order.ID.String(),
	})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	order, err := h.orderService.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sort", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	list, err := h.orderService.ListUserOrders(userID, page, limit, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": list})
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	orderID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	status := models.OrderStatus(req.Status)
	order, err := h.orderService.UpdateOrderStatus(orderID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

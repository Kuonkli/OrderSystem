package main

import (
	"OrderSystem/pkg/database"
	"OrderSystem/pkg/logger"
	"OrderSystem/services/orders/internal/handlers"
	"OrderSystem/services/orders/internal/models"
	"OrderSystem/services/orders/internal/repository"
	"OrderSystem/services/orders/internal/service"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	log := logger.Init("ORDERS")

	dbConfig := database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Info("Connected to database")

	gormDB := db.GetDB()
	err = gormDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		log.Fatal("Failed to create uuid extension:", err)
	}

	err = gormDB.AutoMigrate(&models.Order{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	orderRepo := repository.NewOrderRepository(gormDB, log)
	orderService := service.NewOrdersService(orderRepo, log)
	orderHandler := handlers.NewOrderHandler(orderService, log)

	r := gin.Default()

	// Защищённые роуты
	protected := r.Group("/")
	protected.Use(func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	})

	protected.POST("/orders", orderHandler.CreateOrder)
	protected.GET("/orders", orderHandler.ListOrders)
	protected.GET("/orders/:id", orderHandler.GetOrder)
	protected.PUT("/orders/:id/status", orderHandler.UpdateStatus)

	log.Info("Orders service running on :8082")
	err = r.Run(":8082")
	if err != nil {
		log.Fatal("Failed to start orders service:", err)
	}
}

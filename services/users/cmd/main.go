package main

import (
	"OrderSystem/pkg/database"
	"OrderSystem/pkg/logger"
	"OrderSystem/pkg/tokens"
	"OrderSystem/services/users/internal/handlers"
	"OrderSystem/services/users/internal/models"
	"OrderSystem/services/users/internal/repository"
	"OrderSystem/services/users/internal/service"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	log := logger.Init("USERS")

	//err := godotenv.Load("../../.env")
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//	return
	//}

	// Конфигурация базы данных
	dbConfig := database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// Подключение к базе данных
	db, err := database.Connect(dbConfig) //
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Info("Connected to database")

	gormDB := db.GetDB()
	err = gormDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		log.Fatal("Failed to create uuid extension:", err)
	}

	err = gormDB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	userRepo := repository.NewUserRepository(db.GetDB(), log)

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	tokenService := tokens.NewTokenService(jwtSecret, log)

	userService := service.NewUsersService(userRepo, log)
	userHandler := handlers.NewUserHandler(userService, tokenService, log)

	r := gin.Default()
	r.POST("/users/register", userHandler.SignUp)
	r.POST("/users/login", userHandler.Login)
	r.GET("/users/profile", userHandler.GetProfile)
	r.PUT("/users/profile", userHandler.UpdateProfile)

	log.Info("Users service running on :8081")
	err = r.Run(":8081")
	if err != nil {
		log.Fatal("Failed to start users service:", err)
		return
	}

}

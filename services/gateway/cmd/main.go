package main

import (
	"OrderSystem/pkg/logger"
	"OrderSystem/pkg/tokens"
	"OrderSystem/services/gateway/internal/handlers"
	"OrderSystem/services/gateway/internal/middleware"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	log := logger.Init("GATEWAY")

	//err := godotenv.Load("../../.env")
	//if err != nil {
	//	log.Fatal("Error loading .env file", err)
	//}

	// Users-service URL из env
	usersServiceURL := os.Getenv("USERS_SERVICE_URL")
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	tokenService := tokens.NewTokenService(jwtKey, log)
	usersProxy := handlers.NewUsersProxy(usersServiceURL, log)

	router := gin.Default()

	api := router.Group("/api")
	{
		public := api.Group("/public")
		{
			public.POST("/signup", usersProxy.ProxyTo("/users/register"))
			public.POST("/login", usersProxy.ProxyTo("/users/login"))
		}
		protected := api.Group("/protected", middleware.AuthMiddleware(tokenService))
		{
			protected.POST("/logout", usersProxy.ProxyTo("/users/logout"))
			protected.GET("/user/profile", usersProxy.ProxyTo("/users/profile"))
			protected.PUT("/user/profile", usersProxy.ProxyTo("/users/profile"))
		}
	}

	// Для защищенных роутов (прокси к другим сервисам) используй миддлвэр
	// Например, router.Group("/api/users").Use(middleware.AuthMiddleware(jwtKey))
	// Затем проксирование к users-service

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	} // Порт gateway
	log.Info("Gateway running on :8080")
}

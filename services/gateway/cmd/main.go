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
	ordersServiceURL := os.Getenv("ORDERS_SERVICE_URL")
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	tokenService := tokens.NewTokenService(jwtKey, log)
	usersProxy := handlers.NewUsersProxy(usersServiceURL, log)
	ordersProxy := handlers.NewOrdersProxy(ordersServiceURL, log)

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
			users := protected.Group("/user")
			{
				users.GET("/profile", usersProxy.ProxyTo("/users/profile"))
				users.PUT("/profile", usersProxy.ProxyTo("/users/profile"))
			}
			order := protected.Group("/order")
			{
				order.POST("/add", ordersProxy.ProxyTo("/orders"))
				order.GET("/list", ordersProxy.ProxyTo("/orders"))
				order.GET("/:id", func(c *gin.Context) {
					ordersProxy.ProxyTo("/orders/" + c.Param("id"))(c)
				})
				order.PUT("/:id/status", func(c *gin.Context) {
					ordersProxy.ProxyTo("/orders/" + c.Param("id") + "/status")(c)
				})
			}

		}
	}

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	} // Порт gateway
	log.Info("Gateway running on :8080")
}

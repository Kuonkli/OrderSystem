package middleware

import (
	"OrderSystem/pkg/tokens"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(tokenService *tokens.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessStr, err := c.Cookie("access_token")
		if err != nil || accessStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing access token"})
			c.Abort()
			return
		}
		refreshStr, err := c.Cookie("refresh_token")
		if err != nil || refreshStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
			c.Abort()
			return
		}

		// Парсим access токен
		claims, err := tokenService.Parse(accessStr)
		if err == nil {
			// Access токен ВАЛИДЕН - используем его
			tokenService.Log.Info("Access token user_id: " + claims.UserID)
			c.Set("user_id", claims.UserID)
			c.Next() // Передаем управление следующему обработчику
			return   // ВАЖНО: завершаем выполнение middleware
		}
		tokenService.Log.Warn("Invalid access token ", err.Error())
		tokenService.Log.Info("Access token expired, refreshing...")
		newAccess, err := tokenService.RefreshAccess(refreshStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			c.Abort()
			return
		}

		c.SetCookie("access_token", newAccess, 900, "/", "", false, true)

		claims, err = tokenService.Parse(newAccess)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		tokenService.Log.Info("Refreshed token user_id: " + claims.UserID)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

package handlers

import (
	"OrderSystem/pkg/dto"
	"OrderSystem/pkg/logger"
	"OrderSystem/pkg/tokens"
	"OrderSystem/services/users/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserHandler struct {
	tokenService *tokens.TokenService
	userService  *service.UsersService
	logger       *logger.Logger
}

func NewUserHandler(userService *service.UsersService, tokenService *tokens.TokenService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService:  userService,
		tokenService: tokenService,
		logger:       logger,
	}
}

// SignUp создает нового пользователя
func (h *UserHandler) SignUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithField("request_id", c.GetString("request_id")).Warnf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request data",
			},
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("Failed to hash password", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	req.Password = string(hash)

	user, err := h.userService.CreateUser(req)
	if err != nil {
		h.logger.WithField("request_id", c.GetString("request_id")).Errorf("Registration failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "REGISTRATION_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	h.logger.Info("Registration succeeded user UD: ", user.ID.String())

	accessTokenString, err := h.tokenService.GenerateAccess(user.ID.String())
	if err != nil {
		h.logger.WithField("request_id", c.GetString("request_id")).Errorf("Failed to generate access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
	}
	refreshTokenString, err := h.tokenService.GenerateRefresh(user.ID.String())
	if err != nil {
		h.logger.WithField("request_id", c.GetString("request_id")).Errorf("Failed to generate refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
	}

	c.Header("Access-Token", accessTokenString)
	c.Header("Refresh-Token", refreshTokenString)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    user,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid login request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		h.logger.Error("User not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		h.logger.Warn("Invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	h.logger.Info("Login succeeded user UD: ", user.ID.String())

	accessTokenString, err := h.tokenService.GenerateAccess(user.ID.String())
	if err != nil {
		h.logger.WithField("request_id", c.GetString("request_id")).Errorf("Failed to generate access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
	}
	refreshTokenString, err := h.tokenService.GenerateRefresh(user.ID.String())
	if err != nil {
		h.logger.WithField("request_id", c.GetString("request_id")).Errorf("Failed to generate refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
	}

	c.Header("Access-Token", accessTokenString)
	c.Header("Refresh-Token", refreshTokenString)

	// Скрываем пароль
	response := &dto.LoginResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Password: user.PasswordHash,
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	h.logger.Info("Get user profile: ", userID)
	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.userService.UpdateProfile(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

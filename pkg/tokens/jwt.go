package tokens

import (
	"OrderSystem/pkg/logger"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type TokenService struct {
	secret []byte
	Log    *logger.Logger
}

func NewTokenService(secret []byte, log *logger.Logger) *TokenService {
	return &TokenService{secret: secret, Log: log}
}

func (s *TokenService) GenerateAccess(userID string) (string, error) {
	s.Log.Info("User ID in generating ", userID)
	accessTokenExpirationTime := time.Now().Add(15 * time.Minute)
	accessClaims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpirationTime.Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.secret)
	if err != nil {
		return "", err
	}
	return accessTokenString, nil
}

func (s *TokenService) GenerateRefresh(userID string) (string, error) {
	refreshTokenExpirationTime := time.Now().Add(72 * time.Hour)
	refreshClaims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpirationTime.Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.secret)
	if err != nil {
		return "", err
	}
	return refreshTokenString, nil
}

func (s *TokenService) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})

	s.Log.Info("Token Parse Success ", token)

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	s.Log.Info("Claims ", claims.UserID)
	return claims, nil
}

func (s *TokenService) RefreshAccess(refreshToken string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	return s.GenerateAccess(claims.UserID)
}

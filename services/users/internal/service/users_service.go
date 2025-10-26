package service

import (
	"OrderSystem/pkg/dto"
	"OrderSystem/pkg/logger"
	"OrderSystem/services/users/internal/models"
	"OrderSystem/services/users/internal/repository"
	"github.com/lib/pq"
)

type UsersService struct {
	userRepo repository.UserRepository
	log      *logger.Logger
}

func NewUsersService(userRepo repository.UserRepository, log *logger.Logger) *UsersService {
	return &UsersService{userRepo: userRepo, log: log}
}

func (s *UsersService) CreateUser(payload dto.SignUpRequest) (user *models.User, err error) {
	userData := &models.User{
		FirstName:    payload.Name,
		LastName:     payload.LastName,
		Email:        payload.Email,
		PasswordHash: payload.Password,
		Roles:        pq.StringArray{"user"},
	}
	user, err = s.userRepo.Create(userData)
	if err != nil {
		s.log.Error("Failed to create user", err)
		return nil, err
	}
	return user, nil
}

func (s *UsersService) GetUserByEmail(email string) (*models.User, error) {
	s.log.Info("Fetching user by email", "email", email)
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		s.log.Warn("User not found", "email", email)
		return nil, err
	}
	return user, nil
}

func (s *UsersService) GetProfile(userID string) (*dto.UserProfile, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return &dto.UserProfile{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Roles:     user.Roles,
	}, nil
}

func (s *UsersService) UpdateProfile(userID string, payload dto.UpdateProfileRequest) (*dto.UserProfile, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if payload.FirstName != "" {
		user.FirstName = payload.FirstName
	}
	if payload.LastName != "" {
		user.LastName = payload.LastName
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &dto.UserProfile{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Roles:     user.Roles,
	}, nil
}

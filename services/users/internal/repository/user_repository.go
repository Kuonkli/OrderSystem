package repository

import (
	"OrderSystem/pkg/logger"
	"OrderSystem/services/users/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	List(limit, offset int) ([]*models.User, int64, error)
}

type userRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewUserRepository(db *gorm.DB, log *logger.Logger) UserRepository {
	return &userRepository{db: db, log: log}
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	result := r.db.Create(user)
	if result.Error != nil {
		r.log.Errorf("Failed to create user: %v", result.Error)
		return nil, result.Error
	}
	return user, nil
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		r.log.Errorf("Failed to update user: %v", result.Error)
		return result.Error
	}
	return nil
}

func (r *userRepository) List(limit, offset int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	// Получаем общее количество
	r.db.Model(&models.User{}).Count(&total)

	// Получаем пользователей с пагинацией
	result := r.db.Offset(offset).Limit(limit).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

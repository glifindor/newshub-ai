package service

import (
	"context"
	"errors"

	"auth-service/internal/models"
	"auth-service/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetProfile получает профиль пользователя
func (s *UserService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// UpdateProfile обновляет профиль пользователя
func (s *UserService) UpdateProfile(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}

// GetUserByID получает пользователя по ID
func (s *UserService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// ListUsers получает список пользователей с пагинацией
func (s *UserService) ListUsers(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.List(ctx, pageSize, offset)
}

// UpdateUserRole обновляет роль пользователя
func (s *UserService) UpdateUserRole(ctx context.Context, userID, role string) error {
	// Проверка валидности роли
	validRoles := map[string]bool{
		"admin":     true,
		"editor":    true,
		"moderator": true,
		"user":      true,
	}

	if !validRoles[role] {
		return errors.New("invalid role")
	}

	return s.userRepo.UpdateRole(ctx, userID, role)
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	return s.userRepo.Delete(ctx, userID)
}

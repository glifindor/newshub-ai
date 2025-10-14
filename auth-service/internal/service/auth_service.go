package service

import (
	"context"
	"errors"

	"auth-service/internal/models"
	"auth-service/internal/repository"
	"auth-service/pkg/hash"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService struct {
	userRepo     *repository.UserRepository
	sessionRepo  *repository.SessionRepository
	tokenService *TokenService
}

func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository, tokenService *TokenService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		tokenService: tokenService,
	}
}

func (s *AuthService) Register(ctx context.Context, user *models.User) error {
	// Проверка существования пользователя
	existing, _ := s.userRepo.GetByEmail(ctx, user.Email)
	if existing != nil {
		return ErrUserAlreadyExists
	}

	// Хеширование пароля
	hashedPassword, err := hash.HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword

	// Установка роли по умолчанию
	if user.Role == "" {
		user.Role = "user"
	}

	// Сохранение пользователя
	return s.userRepo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.TokenPair, error) {
	// Получение пользователя
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Проверка пароля
	if !hash.CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Генерация токенов
	tokens, err := s.tokenService.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	// Сохранение refresh токена в Redis
	if err := s.sessionRepo.SaveRefreshToken(ctx, user.ID, tokens.RefreshToken, tokens.RefreshExpiresIn); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*models.TokenClaims, error) {
	return s.tokenService.ValidateAccessToken(token)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	// Валидация refresh токена
	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Проверка существования токена в Redis
	exists, err := s.sessionRepo.ValidateRefreshToken(ctx, claims.UserID, refreshToken)
	if err != nil || !exists {
		return nil, errors.New("invalid refresh token")
	}

	// Получение пользователя
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Генерация новой пары токенов
	tokens, err := s.tokenService.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	// Удаление старого refresh токена
	s.sessionRepo.DeleteRefreshToken(ctx, claims.UserID, refreshToken)

	// Сохранение нового refresh токена
	if err := s.sessionRepo.SaveRefreshToken(ctx, user.ID, tokens.RefreshToken, tokens.RefreshExpiresIn); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	return s.sessionRepo.DeleteRefreshToken(ctx, claims.UserID, refreshToken)
}

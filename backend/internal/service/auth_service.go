package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/helium-mmm/Publishier/internal/auth"
	"github.com/helium-mmm/Publishier/internal/domain"
	"github.com/helium-mmm/Publishier/internal/repository"
	"github.com/helium-mmm/Publishier/internal/validation"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
	tokens   *auth.TokenManager
}

func NewAuthService(userRepo repository.UserRepository, tokens *auth.TokenManager) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		tokens:   tokens,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return "", domain.ErrInvalidCredentials
	}
	if !validation.IsValidEmail(email) {
		return "", domain.ErrInvalidEmail
	}

	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return "", domain.ErrUserAlreadyExists
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return "", err
	}

	return s.tokens.Generate(user.ID)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	return s.tokens.Generate(user.ID)
}

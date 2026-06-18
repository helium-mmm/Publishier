package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/helium-mmm/Publishier/internal/crypto"
	"github.com/helium-mmm/Publishier/internal/domain"
	"github.com/helium-mmm/Publishier/internal/repository"
)

type AccountService struct {
	accountRepo repository.SocialAccountRepository
	encryptor   *crypto.Encryptor
}

func NewAccountService(accountRepo repository.SocialAccountRepository, encryptor *crypto.Encryptor) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		encryptor:   encryptor,
	}
}

func (s *AccountService) ConnectTelegram(ctx context.Context, userID uuid.UUID, botToken, chatID string) error {
	botToken = strings.TrimSpace(botToken)
	chatID = strings.TrimSpace(chatID)
	if botToken == "" || chatID == "" {
		return domain.ErrInvalidCredentials
	}

	encrypted, err := s.encryptor.Encrypt(botToken)
	if err != nil {
		return err
	}

	existing, err := s.accountRepo.GetByUserAndPlatform(ctx, userID, domain.Telegram)
	if err != nil && !errors.Is(err, domain.ErrAccountNotFound) {
		return err
	}

	if existing != nil {
		existing.EncryptedToken = encrypted
		existing.ChatID = chatID
		return s.accountRepo.Update(ctx, existing)
	}

	account := domain.SocialAccount{
		ID:             uuid.New(),
		UserID:         userID,
		Platform:       domain.Telegram,
		EncryptedToken: encrypted,
		ChatID:         chatID,
	}

	return s.accountRepo.Create(ctx, &account)
}

func (s *AccountService) DecryptToken(account domain.SocialAccount) (string, error) {
	return s.encryptor.Decrypt(account.EncryptedToken)
}

type TelegramStatus struct {
	Connected bool
	ChatID    string
	Platform  domain.Platform
}

func (s *AccountService) GetTelegramStatus(ctx context.Context, userID uuid.UUID) (*TelegramStatus, error) {
	account, err := s.accountRepo.GetByUserAndPlatform(ctx, userID, domain.Telegram)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return &TelegramStatus{Connected: false}, nil
		}
		return nil, err
	}

	return &TelegramStatus{
		Connected: true,
		ChatID:    account.ChatID,
		Platform:  account.Platform,
	}, nil
}

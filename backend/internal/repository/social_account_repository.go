package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/helium-mmm/Publishier/internal/domain"
)

type SocialAccountRepository interface {
	Create(ctx context.Context, account *domain.SocialAccount) error
	Update(ctx context.Context, account *domain.SocialAccount) error
	GetByUserAndPlatform(ctx context.Context, userID uuid.UUID, platform domain.Platform) (*domain.SocialAccount, error)
}

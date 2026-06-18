package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/helium-mmm/Publishier/internal/domain"
	"github.com/helium-mmm/Publishier/internal/repository/postgres/models"
	"gorm.io/gorm"
)

type SocialAccountRepository struct {
	db *gorm.DB
}

func NewSocialAccountRepository(db *gorm.DB) *SocialAccountRepository {
	return &SocialAccountRepository{db: db}
}

func (r *SocialAccountRepository) Create(ctx context.Context, account *domain.SocialAccount) error {
	m := accountToModel(*account)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	*account = accountToDomain(m)
	return nil
}

func (r *SocialAccountRepository) Update(ctx context.Context, account *domain.SocialAccount) error {
	m := accountToModel(*account)
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *SocialAccountRepository) GetByUserAndPlatform(ctx context.Context, userID uuid.UUID, platform domain.Platform) (*domain.SocialAccount, error) {
	var m models.SocialAccount
	err := r.db.WithContext(ctx).First(&m, "user_id = ? AND platform = ?", userID, platform).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}
	account := accountToDomain(m)
	return &account, nil
}

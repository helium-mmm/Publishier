package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/helium-mmm/Publishier/internal/domain"
	"github.com/helium-mmm/Publishier/internal/repository/postgres/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	m := userToModel(*user)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	*user = userToDomain(m)
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m models.User
	err := r.db.WithContext(ctx).First(&m, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	user := userToDomain(m)
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var m models.User
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	user := userToDomain(m)
	return &user, nil
}

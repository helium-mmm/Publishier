package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/helium-mmm/Publishier/internal/domain"
	"github.com/helium-mmm/Publishier/internal/repository/postgres/models"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	m := postToModel(*post)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	*post = postToDomain(m)
	return nil
}

func (r *PostRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	var m models.Post
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrPostNotFound
	}
	if err != nil {
		return nil, err
	}
	post := postToDomain(m)
	return &post, nil
}

func (r *PostRepository) GetByIDAndUser(ctx context.Context, id, userID uuid.UUID) (*domain.Post, error) {
	var m models.Post
	err := r.db.WithContext(ctx).First(&m, "id = ? AND user_id = ?", id, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrPostNotFound
	}
	if err != nil {
		return nil, err
	}
	post := postToDomain(m)
	return &post, nil
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	m := postToModel(*post)
	return r.db.WithContext(ctx).Save(&m).Error
}

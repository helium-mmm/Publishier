package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/helium-mmm/Publishier/internal/domain"
)

type PostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Post, error)
	GetByIDAndUser(ctx context.Context, id, userID uuid.UUID) (*domain.Post, error)
	Update(ctx context.Context, post *domain.Post) error
}

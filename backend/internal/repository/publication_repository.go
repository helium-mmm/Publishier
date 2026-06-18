package repository

import (
	"context"

	"github.com/helium-mmm/Publishier/internal/domain"
)

type PublicationRepository interface {
	Create(ctx context.Context, pub *domain.Publication) error
	Update(ctx context.Context, pub *domain.Publication) error
}

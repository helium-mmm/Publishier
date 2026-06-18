package postgres

import (
	"context"

	"github.com/helium-mmm/Publishier/internal/domain"
	"gorm.io/gorm"
)

type PublicationRepository struct {
	db *gorm.DB
}

func NewPublicationRepository(db *gorm.DB) *PublicationRepository {
	return &PublicationRepository{db: db}
}

func (r *PublicationRepository) Create(ctx context.Context, pub *domain.Publication) error {
	m := publicationToModel(*pub)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	*pub = publicationToDomain(m)
	return nil
}

func (r *PublicationRepository) Update(ctx context.Context, pub *domain.Publication) error {
	m := publicationToModel(*pub)
	return r.db.WithContext(ctx).Save(&m).Error
}

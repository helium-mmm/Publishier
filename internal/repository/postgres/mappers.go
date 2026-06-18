package postgres

import (
	"github.com/helium-mmm/Publishier/internal/domain"
	"github.com/helium-mmm/Publishier/internal/repository/postgres/models"
)

func toDomain(m models.Post) domain.Post {
	return domain.Post{
		ID: m.ID,
		Content: m.Content,
		Status: domain.PostStatus(m.Status),
		CreatedAt: m.CreatedAt,
		PublishedAt: m.PublishedAt,
	}
}

func toModel(d domain.Post) models.Post {
	return models.Post{
		ID: d.ID,
		Content: d.Content,
		Status: string(d.Status),
		CreatedAt: d.CreatedAt,
		PublishedAt: d.PublishedAt,
	}
}
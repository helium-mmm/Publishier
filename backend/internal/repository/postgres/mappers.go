package postgres

import (
	"github.com/helium-mmm/Publishier/internal/domain"
	"github.com/helium-mmm/Publishier/internal/repository/postgres/models"
)

func userToDomain(m models.User) domain.User {
	return domain.User{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
	}
}

func userToModel(d domain.User) models.User {
	return models.User{
		ID:           d.ID,
		Email:        d.Email,
		PasswordHash: d.PasswordHash,
	}
}

func postToDomain(m models.Post) domain.Post {
	return domain.Post{
		ID:          m.ID,
		UserID:      m.UserID,
		Content:     m.Content,
		Status:      domain.PostStatus(m.Status),
		CreatedAt:   m.CreatedAt,
		PublishedAt: m.PublishedAt,
	}
}

func postToModel(d domain.Post) models.Post {
	return models.Post{
		ID:          d.ID,
		UserID:      d.UserID,
		Content:     d.Content,
		Status:      string(d.Status),
		CreatedAt:   d.CreatedAt,
		PublishedAt: d.PublishedAt,
	}
}

func accountToDomain(m models.SocialAccount) domain.SocialAccount {
	return domain.SocialAccount{
		ID:             m.ID,
		UserID:         m.UserID,
		Platform:       domain.Platform(m.Platform),
		EncryptedToken: m.EncryptedToken,
		ChatID:         m.ChatID,
	}
}

func accountToModel(d domain.SocialAccount) models.SocialAccount {
	return models.SocialAccount{
		ID:             d.ID,
		UserID:         d.UserID,
		Platform:       string(d.Platform),
		EncryptedToken: d.EncryptedToken,
		ChatID:         d.ChatID,
	}
}

func publicationToDomain(m models.Publication) domain.Publication {
	return domain.Publication{
		ID:         m.ID,
		PostID:     m.PostID,
		UserID:     m.UserID,
		Platform:   domain.Platform(m.Platform),
		Status:     domain.PublicationStatus(m.Status),
		ExternalID: m.ExternalID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func publicationToModel(d domain.Publication) models.Publication {
	return models.Publication{
		ID:         d.ID,
		PostID:     d.PostID,
		UserID:     d.UserID,
		Platform:   string(d.Platform),
		Status:     string(d.Status),
		ExternalID: d.ExternalID,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}

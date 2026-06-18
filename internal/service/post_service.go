package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/helium-mmm/Publishier/internal/domain"
	"github.com/helium-mmm/Publishier/internal/publishier"
	"github.com/helium-mmm/Publishier/internal/repository"
)

type PostService struct {
	repo repository.PostRepository
	accountRepo 
	publicationRepo
	telegramPublishier publishier.Publishier
}

func NewPostService(repo repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Create(ctx context.Context, content string) error {
	post := domain.Post{
		ID: uuid.New(),
		Content: content,
		Status: domain.Draft,
	}
	return s.repo.Create(ctx, &post)
}
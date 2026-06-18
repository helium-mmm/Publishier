package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/helium-mmm/Publishier/internal/domain"
	"github.com/helium-mmm/Publishier/internal/publishier"
	"github.com/helium-mmm/Publishier/internal/repository"
)

type PostService struct {
	repo               repository.PostRepository
	accountRepo        repository.SocialAccountRepository
	publicationRepo    repository.PublicationRepository
	accountService     *AccountService
	telegramPublishier publishier.Publishier
}

func NewPostService(
	repo repository.PostRepository,
	accountRepo repository.SocialAccountRepository,
	publicationRepo repository.PublicationRepository,
	accountService *AccountService,
	telegramPublishier publishier.Publishier,
) *PostService {
	return &PostService{
		repo:               repo,
		accountRepo:        accountRepo,
		publicationRepo:    publicationRepo,
		accountService:     accountService,
		telegramPublishier: telegramPublishier,
	}
}

func (s *PostService) Create(ctx context.Context, userID uuid.UUID, content string) (*domain.Post, error) {
	now := time.Now()
	post := domain.Post{
		ID:        uuid.New(),
		UserID:    userID,
		Content:   content,
		Status:    domain.Draft,
		CreatedAt: now,
	}

	if err := s.repo.Create(ctx, &post); err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *PostService) Get(ctx context.Context, userID, postID uuid.UUID) (*domain.Post, error) {
	return s.repo.GetByIDAndUser(ctx, postID, userID)
}

func (s *PostService) Publish(ctx context.Context, userID, postID uuid.UUID) (*domain.Post, error) {
	post, err := s.repo.GetByIDAndUser(ctx, postID, userID)
	if err != nil {
		return nil, err
	}

	if post.Status != domain.Draft {
		return nil, domain.ErrInvalidStatus
	}

	account, err := s.accountRepo.GetByUserAndPlatform(ctx, userID, domain.Telegram)
	if err != nil {
		return nil, err
	}

	token, err := s.accountService.DecryptToken(*account)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	publication := domain.Publication{
		ID:        uuid.New(),
		PostID:    post.ID,
		UserID:    userID,
		Platform:  domain.Telegram,
		Status:    domain.PublicationPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.publicationRepo.Create(ctx, &publication); err != nil {
		return nil, err
	}

	externalID, err := s.telegramPublishier.Publish(ctx, *post, *account, token)
	if err != nil {
		publication.Status = domain.PublicationFailed
		publication.UpdatedAt = time.Now()
		_ = s.publicationRepo.Update(ctx, &publication)

		post.Status = domain.Failed
		_ = s.repo.Update(ctx, post)
		return nil, domain.ErrPublicationFailed
	}

	publication.Status = domain.PublicationSuccess
	publication.ExternalID = &externalID
	publication.UpdatedAt = time.Now()
	if err := s.publicationRepo.Update(ctx, &publication); err != nil {
		return nil, err
	}

	post.Status = domain.Published
	post.PublishedAt = &now
	if err := s.repo.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

package domain

import (
	"time"

	"github.com/google/uuid"
)

type PostStatus string

const (
	Draft PostStatus = "DRAFT"
	Published PostStatus = "PUBLISHED"
	Failed PostStatus = "FAILED"
)

type Post struct {
	ID uuid.UUID 
	UserID uuid.UUID

	Content   string

	Status PostStatus
	
	CreatedAt time.Time
	PublishedAt *time.Time
}


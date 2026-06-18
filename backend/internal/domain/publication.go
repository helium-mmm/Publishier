package domain

import (
	"time"

	"github.com/google/uuid"
)

type PublicationStatus string

const (
	PublicationPending  PublicationStatus = "PENDING"
	PublicationSuccess PublicationStatus = "SUCCESS"
	PublicationFailed PublicationStatus = "FAILED"
)

type Publication struct {
	ID uuid.UUID

	PostID uuid.UUID
	UserID uuid.UUID

	Platform Platform
	Status PublicationStatus

	ExternalID *string //id поста в соц сети

	CreatedAt time.Time
	UpdatedAt time.Time
}
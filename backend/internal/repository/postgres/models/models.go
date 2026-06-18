package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
}

type Post struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Content     string    `gorm:"not null"`
	Status      string    `gorm:"not null"`
	CreatedAt   time.Time
	PublishedAt *time.Time
}

type SocialAccount struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_platform"`
	Platform       string    `gorm:"not null;uniqueIndex:idx_user_platform"`
	EncryptedToken string    `gorm:"not null"`
	ChatID         string    `gorm:"not null"`
	CreatedAt      time.Time
}

type Publication struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	PostID     uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Platform   string    `gorm:"not null"`
	Status     string    `gorm:"not null"`
	ExternalID *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

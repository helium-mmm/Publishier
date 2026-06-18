package domain

import "github.com/google/uuid"

type Platform string 

const (
	Telegram Platform = "TELEGRAM"
)

type SocialAccount struct {
	ID uuid.UUID
	UserID uuid.UUID

	Platform Platform
	EncryptedToken string
	ChatID string
}
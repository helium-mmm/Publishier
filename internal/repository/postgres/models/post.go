package models

import (
	"time"
	"github.com/google/uuid"
)

type Post struct {
	ID uuid.UUID `gorm: "type:uuid;primaryKey"`
	Content   string
	Status string
	CreatedAt time.Time
	PublishedAt *time.Time
}
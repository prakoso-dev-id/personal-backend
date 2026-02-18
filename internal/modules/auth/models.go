package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/prakoso-id/personal-backend/internal/modules/profiles"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string    `gorm:"type:varchar(255);unique;not null"`
	PasswordHash string            `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Profile      *profiles.Profile `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}

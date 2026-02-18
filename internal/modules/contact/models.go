package contact

import (
	"time"

	"github.com/google/uuid"
)

type ContactMessage struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(255);not null"`
	Subject   string    `gorm:"type:varchar(255)"`
	Message   string    `gorm:"type:text;not null"`
	Status    string    `gorm:"type:varchar(50);default:'unread'"`
	CreatedAt time.Time
}

func (ContactMessage) TableName() string {
	return "contact_messages"
}

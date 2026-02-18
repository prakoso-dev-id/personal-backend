package posts

import (
	"time"

	"github.com/google/uuid"
	"github.com/prakoso-id/personal-backend/internal/modules/images"
)

type Post struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title           string          `gorm:"type:varchar(255);not null"`
	Slug            string          `gorm:"type:varchar(255);unique;not null"`
	ContentMarkdown string          `gorm:"type:text"`
	Summary         string          `gorm:"type:text"`
	IsPublished     bool            `gorm:"default:false"`
	PublishedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Tags            []*Tag          `gorm:"many2many:post_tags;"`
	Images          []images.Image  `gorm:"polymorphic:Entity;polymorphicValue:post"`
}

type Tag struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name string    `gorm:"type:varchar(100);unique;not null"`
	Slug string    `gorm:"type:varchar(100);unique;not null"`
}

func (Post) TableName() string {
	return "posts"
}

func (Tag) TableName() string {
	return "tags"
}

package projects

import (
	"time"

	"github.com/google/uuid"
	"github.com/prakoso-id/personal-backend/internal/modules/images"
	"github.com/prakoso-id/personal-backend/internal/modules/skills"
)

type Project struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title           string          `gorm:"type:varchar(255);not null"`
	Slug            string          `gorm:"type:varchar(255);unique;not null"`
	Description     string          `gorm:"type:text"`
	ContentMarkdown string          `gorm:"type:text"`
	DemoURL         string          `gorm:"type:varchar(255)"`
	RepoURL         string          `gorm:"type:varchar(255)"`
	StartDate       *time.Time      `gorm:"type:date"`
	EndDate         *time.Time      `gorm:"type:date"`
	IsFeatured      bool            `gorm:"default:false"`
	ExperienceID    *uuid.UUID      `gorm:"type:uuid;default:null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Skills          []*skills.Skill `gorm:"many2many:project_skills;"`
	Images          []images.Image  `gorm:"polymorphic:Entity;polymorphicValue:project"`
}

func (Project) TableName() string {
	return "projects"
}

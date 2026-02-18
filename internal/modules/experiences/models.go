package experiences

import (
	"time"

	"github.com/google/uuid"
	"github.com/prakoso-id/personal-backend/internal/modules/projects"
	"gorm.io/gorm"
)

type Experience struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Company      string    `gorm:"type:varchar(255);not null"`
	Position     string    `gorm:"type:varchar(255);not null"`
	Description  string    `gorm:"type:text"`
	StartDate    time.Time `gorm:"not null"`
	EndDate      *time.Time
	IsCurrent    bool               `gorm:"default:false"`
	ProfileID    uuid.UUID          `gorm:"type:uuid;not null"`
	Projects     []projects.Project `gorm:"foreignKey:ExperienceID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Experience) TableName() string {
	return "experiences"
}

// Hooks

func (e *Experience) BeforeCreate(tx *gorm.DB) (err error) {
	if e.IsCurrent {
		e.EndDate = nil
	}
	return
}

func (e *Experience) BeforeUpdate(tx *gorm.DB) (err error) {
	if e.IsCurrent {
		e.EndDate = nil
	}
	return
}

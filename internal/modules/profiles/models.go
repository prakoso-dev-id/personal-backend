package profiles

import (
	"time"

	"github.com/google/uuid"
	"github.com/prakoso-id/personal-backend/internal/modules/skills"
)

type Profile struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID       `gorm:"type:uuid;unique;not null"`
	FullName    string          `gorm:"type:varchar(255)"`
	Bio         string          `gorm:"type:text"`
	AvatarURL   string          `gorm:"type:varchar(255)"`
	ResumeURL   string          `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SocialLinks []SocialLink    `gorm:"foreignKey:ProfileID"`
	Experiences []Experience    `gorm:"foreignKey:ProfileID"`
	Skills      []*skills.Skill `gorm:"many2many:profile_skills;"`
}

type SocialLink struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProfileID  uuid.UUID `gorm:"type:uuid;not null"`
	Platform   string    `gorm:"type:varchar(50)"`
	URL        string    `gorm:"type:varchar(255)"`
	OrderIndex int       `gorm:"default:0"`
}

type Experience struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProfileID   uuid.UUID `gorm:"type:uuid;not null"`
	Position    string    `gorm:"type:varchar(255);not null"`
	Company     string    `gorm:"type:varchar(255);not null"`
	Location    string    `gorm:"type:varchar(255)"`
	StartDate   time.Time `gorm:"type:date;not null"`
	EndDate     *time.Time `gorm:"type:date"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Profile) TableName() string {
	return "profiles"
}

func (SocialLink) TableName() string {
	return "social_links"
}

func (Experience) TableName() string {
	return "experiences"
}

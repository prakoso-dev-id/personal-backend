package profiles

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/prakoso-id/personal-backend/internal/modules/skills"
	"gorm.io/gorm"
)

// Package-level base URL for constructing full file URLs at response time.
var baseURL string

// SetBaseURL sets the base URL used to construct full file URLs.
// Should be called once during application initialization.
func SetBaseURL(url string) {
	baseURL = strings.TrimRight(url, "/")
}

type Profile struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID       `gorm:"type:uuid;unique;not null"`
	FullName    string          `gorm:"type:varchar(255)"`
	Bio         string          `gorm:"type:text"`
	AvatarURL   string          `gorm:"type:varchar(512)"`
	ResumeURL   string          `gorm:"type:varchar(512)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SocialLinks []SocialLink    `gorm:"foreignKey:ProfileID"`
	Experiences []Experience    `gorm:"foreignKey:ProfileID"`
	Skills      []*skills.Skill `gorm:"many2many:profile_skills;"`
}

// AfterFind is a GORM hook that prepends the base URL to AvatarURL and ResumeURL
// after loading from the database, so the API response contains the full URL.
func (p *Profile) AfterFind(tx *gorm.DB) error {
	if baseURL != "" {
		if p.AvatarURL != "" && !strings.HasPrefix(p.AvatarURL, "http") {
			p.AvatarURL = baseURL + p.AvatarURL
		}
		if p.ResumeURL != "" && !strings.HasPrefix(p.ResumeURL, "http") {
			p.ResumeURL = baseURL + p.ResumeURL
		}
	}
	return nil
}

// BeforeSave is a GORM hook that strips the base URL from AvatarURL and ResumeURL
// before saving to the database, ensuring only relative paths are stored.
func (p *Profile) BeforeSave(tx *gorm.DB) error {
	if baseURL != "" {
		if strings.HasPrefix(p.AvatarURL, baseURL) {
			p.AvatarURL = strings.TrimPrefix(p.AvatarURL, baseURL)
		}
		if strings.HasPrefix(p.ResumeURL, baseURL) {
			p.ResumeURL = strings.TrimPrefix(p.ResumeURL, baseURL)
		}
	}
	return nil
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

package images

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Package-level base URL for constructing full image URLs at response time.
var baseURL string

// SetBaseURL sets the base URL used to construct full image URLs.
// Should be called once during application initialization.
func SetBaseURL(url string) {
	baseURL = strings.TrimRight(url, "/")
}

type Image struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EntityType  string    `gorm:"type:varchar(50);not null"`
	EntityID    uuid.UUID `gorm:"type:uuid;not null"`
	FileName    string    `gorm:"type:varchar(255);not null"`
	FilePath    string    `gorm:"type:varchar(255);not null"`
	MimeType    string    `gorm:"type:varchar(50)"`
	Size        int64
	AltText     string    `gorm:"type:varchar(255)"`
	IsPrimary   bool      `gorm:"default:false"`
	OrderIndex  int       `gorm:"default:0"`
	CreatedAt   time.Time
}

type ImageUploadResult struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

func (Image) TableName() string {
	return "images"
}

// AfterFind is a GORM hook that prepends the base URL to FilePath
// after loading from the database, so the API response contains the full URL.
func (img *Image) AfterFind(tx *gorm.DB) error {
	if baseURL != "" && img.FilePath != "" && !strings.HasPrefix(img.FilePath, "http") {
		img.FilePath = baseURL + img.FilePath
	}
	return nil
}

// BeforeSave is a GORM hook that strips the base URL from FilePath
// before saving to the database, ensuring only relative paths are stored.
func (img *Image) BeforeSave(tx *gorm.DB) error {
	if baseURL != "" && strings.HasPrefix(img.FilePath, baseURL) {
		img.FilePath = strings.TrimPrefix(img.FilePath, baseURL)
	}
	return nil
}

package skills

import (
	"github.com/google/uuid"
)

type Skill struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name     string    `gorm:"type:varchar(100);unique;not null"`
	Category string    `gorm:"type:varchar(50)"`
	IconURL  string    `gorm:"type:varchar(255)"`
}

func (Skill) TableName() string {
	return "skills"
}

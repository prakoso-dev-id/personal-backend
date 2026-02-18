package database

import (
	"log"

	"github.com/prakoso-id/personal-backend/internal/modules/auth"
	"github.com/prakoso-id/personal-backend/internal/modules/contact"
	"github.com/prakoso-id/personal-backend/internal/modules/experiences"
	"github.com/prakoso-id/personal-backend/internal/modules/images"
	"github.com/prakoso-id/personal-backend/internal/modules/posts"
	"github.com/prakoso-id/personal-backend/internal/modules/profiles"
	"github.com/prakoso-id/personal-backend/internal/modules/projects"
	"github.com/prakoso-id/personal-backend/internal/modules/skills"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	log.Println("Migrating database...")

	err := db.AutoMigrate(
		&auth.User{},
		&profiles.Profile{},
		&profiles.SocialLink{},
		&skills.Skill{},
		&projects.Project{},
		&experiences.Experience{},
		&posts.Post{},
		&posts.Tag{},
		&contact.ContactMessage{},
		&images.Image{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migrated successfully")
}

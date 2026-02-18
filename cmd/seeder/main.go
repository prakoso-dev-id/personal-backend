package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/prakoso-id/personal-backend/internal/config"
	"github.com/prakoso-id/personal-backend/internal/database"
	"github.com/prakoso-id/personal-backend/internal/modules/auth"
	"github.com/prakoso-id/personal-backend/internal/modules/contact"
	"github.com/prakoso-id/personal-backend/internal/modules/images"
	"github.com/prakoso-id/personal-backend/internal/modules/posts"
	"github.com/prakoso-id/personal-backend/internal/modules/profiles"
	"github.com/prakoso-id/personal-backend/internal/modules/projects"
	"github.com/prakoso-id/personal-backend/internal/modules/skills"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect Database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Starting seeding process...")

	// Clean Database
	if err := cleanDB(db); err != nil {
		log.Fatalf("Failed to clean database: %v", err)
	}
	log.Println("Database cleaned.")

	// Seed Data
	user, err := seedUsers(db)
	if err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}
	log.Println("Users seeded.")

	profile, err := seedProfiles(db, user.ID)
	if err != nil {
		log.Fatalf("Failed to seed profiles: %v", err)
	}
	log.Println("Profiles seeded.")

	skillList, err := seedSkills(db)
	if err != nil {
		log.Fatalf("Failed to seed skills: %v", err)
	}
	log.Println("Skills seeded.")

	if err := seedProfileSkills(db, profile.ID, skillList); err != nil {
		log.Fatalf("Failed to seed profile skills: %v", err)
	}
	log.Println("Profile skills seeded.")

	if err := seedExperiences(db, profile.ID); err != nil {
		log.Fatalf("Failed to seed experiences: %v", err)
	}
	log.Println("Experiences seeded.")

	if err := seedSocialLinks(db, profile.ID); err != nil {
		log.Fatalf("Failed to seed social links: %v", err)
	}
	log.Println("Social links seeded.")

	if err := seedProjects(db, skillList); err != nil {
		log.Fatalf("Failed to seed projects: %v", err)
	}
	log.Println("Projects seeded.")

	tags, err := seedTags(db)
	if err != nil {
		log.Fatalf("Failed to seed tags: %v", err)
	}
	log.Println("Tags seeded.")

	if err := seedPosts(db, tags); err != nil {
		log.Fatalf("Failed to seed posts: %v", err)
	}
	log.Println("Posts seeded.")

	if err := seedContactMessages(db); err != nil {
		log.Fatalf("Failed to seed contact messages: %v", err)
	}
	log.Println("Contact messages seeded.")

	log.Println("Seeding completed successfully!")
}

func cleanDB(db *gorm.DB) error {
	// Disable foreign key checks to allow truncation
	if err := db.Exec("TRUNCATE TABLE users, profiles, skills, profile_skills, experiences, social_links, projects, project_skills, tags, posts, post_tags, images, contact_messages RESTART IDENTITY CASCADE").Error; err != nil {
		return err
	}
	return nil
}

func seedUsers(db *gorm.DB) (*auth.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("user1234"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := auth.User{
		Email:        "admin@example.com",
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func seedProfiles(db *gorm.DB, userID uuid.UUID) (*profiles.Profile, error) {
	profile := profiles.Profile{
		UserID:    userID,
		FullName:  "John Doe",
		Bio:       "Full Stack Developer based in Indonesia. Passionate about building scalable web applications.",
		AvatarURL: "https://ui-avatars.com/api/?name=John+Doe",
		ResumeURL: "https://example.com/resume.pdf",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&profile).Error; err != nil {
		return nil, err
	}

	return &profile, nil
}

func seedSkills(db *gorm.DB) ([]*skills.Skill, error) {
	skillsList := []skills.Skill{
		{Name: "Go", Category: "Backend", IconURL: "https://skillicons.dev/icons?i=go"},
		{Name: "PostgreSQL", Category: "Database", IconURL: "https://skillicons.dev/icons?i=postgres"},
		{Name: "React", Category: "Frontend", IconURL: "https://skillicons.dev/icons?i=react"},
		{Name: "TypeScript", Category: "Language", IconURL: "https://skillicons.dev/icons?i=ts"},
		{Name: "Docker", Category: "DevOps", IconURL: "https://skillicons.dev/icons?i=docker"},
		{Name: "AWS", Category: "Cloud", IconURL: "https://skillicons.dev/icons?i=aws"},
	}

	var createdSkills []*skills.Skill
	for _, s := range skillsList {
		if err := db.Create(&s).Error; err != nil {
			return nil, err
		}
		createdSkills = append(createdSkills, &s)
	}

	return createdSkills, nil
}

func seedProfileSkills(db *gorm.DB, profileID uuid.UUID, availableSkills []*skills.Skill) error {
	// Assign first 3 skills to profile
	if len(availableSkills) < 3 {
		return nil
	}
	
	for i := 0; i < 3; i++ {
		// Because GORM many2many can be tricky with just IDs, using raw SQL is safer for simple seeding
		if err := db.Exec("INSERT INTO profile_skills (profile_id, skill_id) VALUES (?, ?)", profileID, availableSkills[i].ID).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedExperiences(db *gorm.DB, profileID uuid.UUID) error {
	now := time.Now()
	exp1 := profiles.Experience{
		ProfileID:   profileID,
		Position:    "Senior Software Engineer",
		Company:     "Tech Corp",
		Location:    "Jakarta, Indonesia",
		StartDate:   now.AddDate(-2, 0, 0),
		EndDate:     nil, // Current job
		Description: "Leading the backend team and architecting scalable solutions.",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	prevDate := now.AddDate(-2, 0, 0)
	exp2 := profiles.Experience{
		ProfileID:   profileID,
		Position:    "Software Engineer",
		Company:     "StartUp Inc",
		Location:    "Bandung, Indonesia",
		StartDate:   now.AddDate(-4, 0, 0),
		EndDate:     &prevDate,
		Description: "Developed key features for the main product.",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := db.Create(&exp1).Error; err != nil {
		return err
	}
	if err := db.Create(&exp2).Error; err != nil {
		return err
	}
	return nil
}

func seedSocialLinks(db *gorm.DB, profileID uuid.UUID) error {
	links := []profiles.SocialLink{
		{ProfileID: profileID, Platform: "github", URL: "https://github.com/johndoe", OrderIndex: 1},
		{ProfileID: profileID, Platform: "linkedin", URL: "https://linkedin.com/in/johndoe", OrderIndex: 2},
		{ProfileID: profileID, Platform: "twitter", URL: "https://twitter.com/johndoe", OrderIndex: 3},
	}

	for _, link := range links {
		if err := db.Create(&link).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedProjects(db *gorm.DB, availableSkills []*skills.Skill) error {
	for i := 0; i < 5; i++ {
		project := projects.Project{
			Title:       faker.Sentence(),
			Slug:        faker.UUIDHyphenated(), // Temporary slug
			Description: faker.Paragraph(),
			ContentMarkdown: "# Project Details\n\n" + faker.Paragraph(),
			DemoURL:     "https://example.com/demo",
			RepoURL:     "https://github.com/johndoe/project",
			IsFeatured:  i < 2, // First 2 are featured
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := db.Create(&project).Error; err != nil {
			return err
		}

		// Add random skills
		if len(availableSkills) > 0 {
			randomSkill := availableSkills[rand.Intn(len(availableSkills))]
			if err := db.Exec("INSERT INTO project_skills (project_id, skill_id) VALUES (?, ?)", project.ID, randomSkill.ID).Error; err != nil {
				return err
			}
		}

		// Add image
		img := images.Image{
			EntityType: "project",
			EntityID:   project.ID,
			FileName:   "project-screenshot.jpg",
			FilePath:   "/uploads/project-" + project.ID.String() + ".jpg", 
			MimeType:   "image/jpeg",
			Size:       102400,
			AltText:    "Project Screenshot",
			IsPrimary:  true,
			CreatedAt:  time.Now(),
		}
		if err := db.Create(&img).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedTags(db *gorm.DB) ([]*posts.Tag, error) {
	tagNames := []string{"Go", "Architecture", "Microservices", "Tutorial", "News"}
	var createdTags []*posts.Tag

	for _, name := range tagNames {
		tag := posts.Tag{
			Name: name,
			Slug: name, // Simplified slug for seeding
		}
		if err := db.Create(&tag).Error; err != nil {
			return nil, err
		}
		createdTags = append(createdTags, &tag)
	}
	return createdTags, nil
}

func seedPosts(db *gorm.DB, availableTags []*posts.Tag) error {
	for i := 0; i < 10; i++ {
		now := time.Now()
		isPublished := i%2 == 0
		var publishedAt *time.Time
		if isPublished {
			publishedAt = &now
		}

		post := posts.Post{
			Title:           faker.Sentence(),
			Slug:            faker.UUIDHyphenated(),
			ContentMarkdown: "# " + faker.Sentence() + "\n\n" + faker.Paragraph(),
			Summary:         faker.Sentence(),
			IsPublished:     isPublished,
			PublishedAt:     publishedAt,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		if err := db.Create(&post).Error; err != nil {
			return err
		}

		// Add random tag
		if len(availableTags) > 0 {
			randomTag := availableTags[rand.Intn(len(availableTags))]
			if err := db.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)", post.ID, randomTag.ID).Error; err != nil {
				return err
			}
		}

		// Add image
		img := images.Image{
			EntityType: "post",
			EntityID:   post.ID,
			FileName:   "post-cover.jpg",
			FilePath:   "/uploads/post-" + post.ID.String() + ".jpg",
			MimeType:   "image/jpeg",
			Size:       204800,
			AltText:    "Post Cover",
			IsPrimary:  true,
			CreatedAt:  now,
		}
		if err := db.Create(&img).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedContactMessages(db *gorm.DB) error {
	for i := 0; i < 5; i++ {
		msg := contact.ContactMessage{
			Name:      faker.Name(),
			Email:     faker.Email(),
			Subject:   faker.Sentence(),
			Message:   faker.Paragraph(),
			Status:    "unread",
			CreatedAt: time.Now(),
		}
		if err := db.Create(&msg).Error; err != nil {
			return err
		}
	}
	return nil
}

package profiles

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetProfile() (*Profile, error)
	GetProfileByUserID(userID uuid.UUID) (*Profile, error)
	Update(profile *Profile) error
	Create(profile *Profile) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetProfile() (*Profile, error) {
	var profile Profile
	// Assuming single profile for personal website, or fetch first one
	err := r.db.Preload("SocialLinks").Preload("Experiences").Preload("Skills").First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *repository) GetProfileByUserID(userID uuid.UUID) (*Profile, error) {
	var profile Profile
	err := r.db.Preload("SocialLinks").Preload("Experiences").Preload("Skills").Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *repository) Create(profile *Profile) error {
	return r.db.Create(profile).Error
}

func (r *repository) Update(profile *Profile) error {
	return r.db.Save(profile).Error
}

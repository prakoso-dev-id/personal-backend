package experiences

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(experience *Experience) error
	Update(experience *Experience) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*Experience, error)
	FindAll(limit, offset int) ([]Experience, error)
	Count() (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(experience *Experience) error {
	return r.db.Create(experience).Error
}

func (r *repository) Update(experience *Experience) error {
	return r.db.Save(experience).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Experience{}, "id = ?", id).Error
}

func (r *repository) FindByID(id uuid.UUID) (*Experience, error) {
	var experience Experience
	err := r.db.Preload("Projects").First(&experience, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &experience, nil
}

func (r *repository) FindAll(limit, offset int) ([]Experience, error) {
	var experiences []Experience
	query := r.db.Order("start_date DESC")
	
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	
	err := query.Preload("Projects").Find(&experiences).Error
	return experiences, err
}

func (r *repository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&Experience{}).Count(&count).Error
	return count, err
}

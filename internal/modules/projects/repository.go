package projects

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(project *Project) error
	Update(project *Project) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*Project, error)
	FindAll(limit, offset int) ([]Project, error)
	Count() (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(project *Project) error {
	return r.db.Create(project).Error
}

func (r *repository) Update(project *Project) error {
	if err := r.db.Model(project).Association("Skills").Replace(project.Skills); err != nil {
		return err
	}
	return r.db.Save(project).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Project{}, "id = ?", id).Error
}

func (r *repository) FindByID(id uuid.UUID) (*Project, error) {
	var project Project
	err := r.db.Preload("Skills").Preload("Images").First(&project, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *repository) FindAll(limit, offset int) ([]Project, error) {
	var projects []Project
	query := r.db.Preload("Skills").Preload("Images").Order("start_date DESC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	err := query.Find(&projects).Error
	return projects, err
}

func (r *repository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&Project{}).Count(&count).Error
	return count, err
}

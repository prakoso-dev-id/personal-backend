package images

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(image *Image) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*Image, error)
	FindByEntity(entityType string, entityID uuid.UUID) ([]Image, error)
	DeleteByEntity(entityType string, entityID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(image *Image) error {
	return r.db.Create(image).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Image{}, "id = ?", id).Error
}

func (r *repository) FindByID(id uuid.UUID) (*Image, error) {
	var image Image
	if err := r.db.First(&image, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &image, nil
}

func (r *repository) FindByEntity(entityType string, entityID uuid.UUID) ([]Image, error) {
	var images []Image
	err := r.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).Order("order_index ASC").Find(&images).Error
	return images, err
}

func (r *repository) DeleteByEntity(entityType string, entityID uuid.UUID) error {
	return r.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).Delete(&Image{}).Error
}

package posts

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(post *Post) error
	Update(post *Post) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*Post, error)
	FindBySlug(slug string) (*Post, error)
	FindAll(publishedOnly bool, limit, offset int) ([]Post, error)
	Count(publishedOnly bool) (int64, error)
	FindOrCreateTag(name, slug string) (*Tag, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindOrCreateTag(name, tagSlug string) (*Tag, error) {
	var tag Tag
	err := r.db.Where("slug = ?", tagSlug).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tag = Tag{Name: name, Slug: tagSlug}
			if err := r.db.Create(&tag).Error; err != nil {
				return nil, err
			}
			return &tag, nil
		}
		return nil, err
	}
	return &tag, nil
}

func (r *repository) Create(post *Post) error {
	return r.db.Create(post).Error
}

func (r *repository) Update(post *Post) error {
	// Update associations (Tags)
	if err := r.db.Model(post).Association("Tags").Replace(post.Tags); err != nil {
		return err
	}
	return r.db.Save(post).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Post{}, "id = ?", id).Error
}

func (r *repository) FindByID(id uuid.UUID) (*Post, error) {
	var post Post
	err := r.db.Preload("Tags").Preload("Images").First(&post, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *repository) FindBySlug(slug string) (*Post, error) {
	var post Post
	err := r.db.Preload("Tags").Preload("Images").First(&post, "slug = ?", slug).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *repository) FindAll(publishedOnly bool, limit, offset int) ([]Post, error) {
	var posts []Post
	query := r.db.Preload("Tags").Preload("Images").Order("created_at DESC")
	if publishedOnly {
		query = query.Where("is_published = ?", true)
	}
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	err := query.Find(&posts).Error
	return posts, err
}

func (r *repository) Count(publishedOnly bool) (int64, error) {
	var count int64
	query := r.db.Model(&Post{})
	if publishedOnly {
		query = query.Where("is_published = ?", true)
	}
	err := query.Count(&count).Error
	return count, err
}

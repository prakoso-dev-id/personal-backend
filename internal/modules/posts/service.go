package posts

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/prakoso-id/personal-backend/internal/modules/images"
	"github.com/prakoso-id/personal-backend/internal/utils/pagination"
)

type Service interface {
	Create(req *CreatePostRequest) (*Post, error)
	Update(id uuid.UUID, req *UpdatePostRequest) (*Post, error)
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*Post, error)
	GetBySlug(slug string) (*Post, error)
	GetAll(public bool) ([]Post, error)
	GetAllAdmin(page, limit int) (*pagination.PaginatedResponse, error)
}

type service struct {
	repo       Repository
	imagesRepo images.Repository
}

func NewService(repo Repository, imagesRepo images.Repository) Service {
	return &service{
		repo:       repo,
		imagesRepo: imagesRepo,
	}
}

type CreatePostRequest struct {
	Title           string                     `json:"title"`
	ContentMarkdown string                     `json:"content_markdown"`
	Summary         string                     `json:"summary"`
	IsPublished     bool                       `json:"is_published"`
	Tags            []string                   `json:"tags"`
	Images          []images.ImageUploadResult `json:"images"`
}

type UpdatePostRequest struct {
	Title           string                     `json:"title"`
	ContentMarkdown string                     `json:"content_markdown"`
	Summary         string                     `json:"summary"`
	IsPublished     bool                       `json:"is_published"`
	Tags            []string                   `json:"tags"`
	Images          []images.ImageUploadResult `json:"images"`
}

func (s *service) Create(req *CreatePostRequest) (*Post, error) {
	post := &Post{
		Title:           req.Title,
		Slug:            slug.Make(req.Title),
		ContentMarkdown: req.ContentMarkdown,
		Summary:         req.Summary,
		IsPublished:     req.IsPublished,
	}

	if req.IsPublished {
		now := time.Now()
		post.PublishedAt = &now
	}

	// Handle tags (find or create)
	var tags []*Tag
	for _, tagName := range req.Tags {
		tag, err := s.repo.FindOrCreateTag(tagName, slug.Make(tagName))
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	post.Tags = tags

	if err := s.repo.Create(post); err != nil {
		return nil, err
	}

	// Handle Images
	for _, imgReq := range req.Images {
		image := &images.Image{
			EntityType: "post",
			EntityID:   post.ID,
			FileName:   imgReq.FileName,
			FilePath:   imgReq.FilePath,
			MimeType:   imgReq.MimeType,
			Size:       imgReq.Size,
			// AltText:    imgReq.AltText, // If we add AltText to ImageUploadResult later
		}
		// We could add error handling here, but maybe just log or ignore?
		// For now, let's try to save.
		_ = s.imagesRepo.Create(image)
	}

	return post, nil
}

func (s *service) Update(id uuid.UUID, req *UpdatePostRequest) (*Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	post.Title = req.Title
	if req.Title != "" {
		post.Slug = slug.Make(req.Title)
	}
	post.ContentMarkdown = req.ContentMarkdown
	post.Summary = req.Summary
	post.IsPublished = req.IsPublished

	if req.IsPublished && post.PublishedAt == nil {
		now := time.Now()
		post.PublishedAt = &now
	}

	// Update tags
	var tags []*Tag
	for _, tagName := range req.Tags {
		tag, err := s.repo.FindOrCreateTag(tagName, slug.Make(tagName))
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	post.Tags = tags

	if err := s.repo.Update(post); err != nil {
		return nil, err
	}

	// Sync Images: delete existing, then insert current set from payload
	if err := s.imagesRepo.DeleteByEntity("post", post.ID); err != nil {
		return nil, err
	}
	for _, imgReq := range req.Images {
		image := &images.Image{
			EntityType: "post",
			EntityID:   post.ID,
			FileName:   imgReq.FileName,
			FilePath:   imgReq.FilePath,
			MimeType:   imgReq.MimeType,
			Size:       imgReq.Size,
		}
		if err := s.imagesRepo.Create(image); err != nil {
			return nil, err
		}
	}

	return post, nil
}

func (s *service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *service) GetByID(id uuid.UUID) (*Post, error) {
	return s.repo.FindByID(id)
}

func (s *service) GetBySlug(slug string) (*Post, error) {
	return s.repo.FindBySlug(slug)
}

func (s *service) GetAll(public bool) ([]Post, error) {
	return s.repo.FindAll(public, 0, 0)
}

func (s *service) GetAllAdmin(page, limit int) (*pagination.PaginatedResponse, error) {
	p := pagination.Pagination{
		Page:  page,
		Limit: limit,
	}

	posts, err := s.repo.FindAll(false, p.Limit, p.Offset())
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(false)
	if err != nil {
		return nil, err
	}

	res := pagination.NewResponse(posts, total, p)
	return &res, nil
}

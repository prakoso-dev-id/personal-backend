package projects

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/prakoso-id/personal-backend/internal/modules/images"
	"github.com/prakoso-id/personal-backend/internal/modules/skills"
	"github.com/prakoso-id/personal-backend/internal/utils/pagination"
)

type Service interface {
	Create(req *CreateProjectRequest) (*Project, error)
	Update(id uuid.UUID, req *UpdateProjectRequest) (*Project, error)
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*Project, error)
	GetAll() ([]Project, error)
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

type CreateProjectRequest struct {
	Title           string                     `json:"title"`
	Description     string                     `json:"description"`
	ContentMarkdown string                     `json:"content_markdown"`
	DemoURL         string                     `json:"demo_url"`
	RepoURL         string                     `json:"repo_url"`
	StartDate       string                     `json:"start_date"` // YYYY-MM-DD
	EndDate         string                     `json:"end_date"`   // YYYY-MM-DD
	IsFeatured      bool                       `json:"is_featured"`
	ExperienceID    *string                    `json:"experience_id"` // UUID or null
	SkillIDs        []string                   `json:"skill_ids"` // UUIDs
	Images          []images.ImageUploadResult `json:"images"`
}

type UpdateProjectRequest struct {
	Title           string                     `json:"title"`
	Description     string                     `json:"description"`
	ContentMarkdown string                     `json:"content_markdown"`
	DemoURL         string                     `json:"demo_url"`
	RepoURL         string                     `json:"repo_url"`
	StartDate       string                     `json:"start_date"`
	EndDate         string                     `json:"end_date"`
	IsFeatured      bool                       `json:"is_featured"`
	ExperienceID    *string                    `json:"experience_id"`
	SkillIDs        []string                   `json:"skill_ids"`
	Images          []images.ImageUploadResult `json:"images"`
}

func parseDate(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}
	return &t
}

func (s *service) Create(req *CreateProjectRequest) (*Project, error) {
	project := &Project{
		Title:           req.Title,
		Slug:            slug.Make(req.Title),
		Description:     req.Description,
		ContentMarkdown: req.ContentMarkdown,
		DemoURL:         req.DemoURL,
		RepoURL:         req.RepoURL,
		StartDate:       parseDate(req.StartDate),
		EndDate:         parseDate(req.EndDate),
		IsFeatured:      req.IsFeatured,
	}

	if req.ExperienceID != nil && *req.ExperienceID != "" {
		id, err := uuid.Parse(*req.ExperienceID)
		if err == nil {
			project.ExperienceID = &id
		}
	}

	var projectSkills []*skills.Skill
	for _, idStr := range req.SkillIDs {
		id, err := uuid.Parse(idStr)
		if err == nil {
			projectSkills = append(projectSkills, &skills.Skill{ID: id})
		}
	}
	project.Skills = projectSkills

	if err := s.repo.Create(project); err != nil {
		return nil, err
	}

	// Handle Images
	for _, imgReq := range req.Images {
		image := &images.Image{
			EntityType: "project",
			EntityID:   project.ID,
			FileName:   imgReq.FileName,
			FilePath:   imgReq.FilePath,
			MimeType:   imgReq.MimeType,
			Size:       imgReq.Size,
		}
		_ = s.imagesRepo.Create(image)
	}

	return project, nil
}

func (s *service) Update(id uuid.UUID, req *UpdateProjectRequest) (*Project, error) {
	project, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	project.Title = req.Title
	if req.Title != "" {
		project.Slug = slug.Make(req.Title)
	}
	project.Description = req.Description
	project.ContentMarkdown = req.ContentMarkdown
	project.DemoURL = req.DemoURL
	project.RepoURL = req.RepoURL
	project.StartDate = parseDate(req.StartDate)
	project.EndDate = parseDate(req.EndDate)
	project.IsFeatured = req.IsFeatured

	if req.ExperienceID != nil {
		if *req.ExperienceID == "" {
			project.ExperienceID = nil
		} else {
			id, err := uuid.Parse(*req.ExperienceID)
			if err == nil {
				project.ExperienceID = &id
			}
		}
	}

	var projectSkills []*skills.Skill
	for _, idStr := range req.SkillIDs {
		id, err := uuid.Parse(idStr)
		if err == nil {
			projectSkills = append(projectSkills, &skills.Skill{ID: id})
		}
	}
	project.Skills = projectSkills

	if err := s.repo.Update(project); err != nil {
		return nil, err
	}

	// Handle New Images
	for _, imgReq := range req.Images {
		image := &images.Image{
			EntityType: "project",
			EntityID:   project.ID,
			FileName:   imgReq.FileName,
			FilePath:   imgReq.FilePath,
			MimeType:   imgReq.MimeType,
			Size:       imgReq.Size,
		}
		_ = s.imagesRepo.Create(image)
	}

	return project, nil
}

func (s *service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *service) GetByID(id uuid.UUID) (*Project, error) {
	return s.repo.FindByID(id)
}

func (s *service) GetAll() ([]Project, error) {
	return s.repo.FindAll(0, 0)
}

func (s *service) GetAllAdmin(page, limit int) (*pagination.PaginatedResponse, error) {
	p := pagination.Pagination{
		Page:  page,
		Limit: limit,
	}

	projects, err := s.repo.FindAll(p.Limit, p.Offset())
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, err
	}

	res := pagination.NewResponse(projects, total, p)
	return &res, nil
}

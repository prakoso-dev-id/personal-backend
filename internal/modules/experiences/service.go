package experiences

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(req *CreateExperienceRequest) (*Experience, error)
	Update(id uuid.UUID, req *UpdateExperienceRequest) (*Experience, error)
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*Experience, error)
	GetAll() ([]Experience, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

type CreateExperienceRequest struct {
	Company      string    `json:"company" binding:"required"`
	Position     string    `json:"position" binding:"required"`
	Description  string    `json:"description"`
	StartDate    string     `json:"start_date" binding:"required"`
	EndDate      string     `json:"end_date"`
	IsCurrent    bool      `json:"is_current"`
	ProfileID    uuid.UUID `json:"-"` // Set by handler
}

type UpdateExperienceRequest struct {
	Company      string    `json:"company"`
	Position     string    `json:"position"`
	Description  string    `json:"description"`
	StartDate    string    `json:"start_date"`
	EndDate      string    `json:"end_date"`
	IsCurrent    bool      `json:"is_current"`
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

func (s *service) Create(req *CreateExperienceRequest) (*Experience, error) {
	startDate := parseDate(req.StartDate)
	if startDate == nil {
		return nil, errors.New("invalid start_date format")
	}

	experience := &Experience{
		Company:      req.Company,
		Position:     req.Position,
		Description:  req.Description,
		StartDate:    *startDate,
		EndDate:      parseDate(req.EndDate),
		IsCurrent:    req.IsCurrent,
		ProfileID:    req.ProfileID,
	}

	if err := s.repo.Create(experience); err != nil {
		return nil, err
	}

	return experience, nil
}

func (s *service) Update(id uuid.UUID, req *UpdateExperienceRequest) (*Experience, error) {
	experience, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if experience == nil {
		return nil, errors.New("experience not found")
	}

	experience.Company = req.Company
	experience.Position = req.Position
	experience.Description = req.Description
	if req.StartDate != "" {
		startDate := parseDate(req.StartDate)
		if startDate != nil {
			experience.StartDate = *startDate
		}
	}
	experience.EndDate = parseDate(req.EndDate)
	experience.IsCurrent = req.IsCurrent


	if err := s.repo.Update(experience); err != nil {
		return nil, err
	}

	return experience, nil
}

func (s *service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *service) GetByID(id uuid.UUID) (*Experience, error) {
	return s.repo.FindByID(id)
}

func (s *service) GetAll() ([]Experience, error) {
	return s.repo.FindAll(0, 0)
}

package profiles

import (
	"github.com/google/uuid"
)

type Service interface {
	GetProfile() (*Profile, error)
	GetProfileByUserID(userID uuid.UUID) (*Profile, error)
	CreateOrUpdateProfile(userID uuid.UUID, req *UpdateProfileRequest) (*Profile, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetProfile() (*Profile, error) {
	return s.repo.GetProfile()
}

func (s *service) GetProfileByUserID(userID uuid.UUID) (*Profile, error) {
	return s.repo.GetProfileByUserID(userID)
}

type UpdateProfileRequest struct {
	FullName  string `json:"full_name"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	ResumeURL string `json:"resume_url"`
}

func (s *service) CreateOrUpdateProfile(userID uuid.UUID, req *UpdateProfileRequest) (*Profile, error) {
	profile, err := s.repo.GetProfileByUserID(userID)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		profile = &Profile{
			UserID: userID,
		}
	}

	profile.FullName = req.FullName
	profile.Bio = req.Bio
	profile.AvatarURL = req.AvatarURL
	profile.ResumeURL = req.ResumeURL

	if profile.ID == uuid.Nil {
		err = s.repo.Create(profile)
	} else {
		err = s.repo.Update(profile)
	}

	if err != nil {
		return nil, err
	}

	return profile, nil
}

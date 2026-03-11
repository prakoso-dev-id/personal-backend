package profiles

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetProfile() (*Profile, error)
	GetProfileByUserID(userID uuid.UUID) (*Profile, error)
	CreateOrUpdateProfile(userID uuid.UUID, req *UpdateProfileRequest) (*Profile, error)
}

type service struct {
	repo    Repository
	storage string
}

func NewService(repo Repository) Service {
	return &service{
		repo:    repo,
		storage: "storage",
	}
}

func (s *service) GetProfile() (*Profile, error) {
	return s.repo.GetProfile()
}

func (s *service) GetProfileByUserID(userID uuid.UUID) (*Profile, error) {
	return s.repo.GetProfileByUserID(userID)
}

type UpdateProfileRequest struct {
	FullName   string                `form:"full_name"`
	Bio        string                `form:"bio"`
	AvatarFile *multipart.FileHeader `form:"avatar"`
	ResumeFile *multipart.FileHeader `form:"resume"`
}

// allowedImageExts defines allowed extensions for avatar images
var allowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
}

// allowedResumeExts defines allowed extensions for resume files
var allowedResumeExts = map[string]bool{
	".pdf": true, ".doc": true, ".docx": true,
}

func (s *service) uploadFile(file *multipart.FileHeader, subDir string, allowedExts map[string]bool, maxSize int64) (string, error) {
	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExts[ext] {
		allowedList := make([]string, 0, len(allowedExts))
		for k := range allowedExts {
			allowedList = append(allowedList, k)
		}
		return "", fmt.Errorf("invalid file type (allowed: %s)", strings.Join(allowedList, ", "))
	}

	// Validate file size
	if file.Size > maxSize {
		return "", fmt.Errorf("file too large (max %dMB)", maxSize/(1024*1024))
	}

	// Create directory structure: storage/{subDir}/
	uploadDir := filepath.Join(s.storage, subDir)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate safe filename
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dstPath := filepath.Join(uploadDir, newFilename)

	// Save file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// Public path for serving: /media/{subDir}/{filename}
	publicPath := fmt.Sprintf("/media/%s/%s", subDir, newFilename)

	// Build full URL for the response
	fullURL := publicPath
	if baseURL != "" {
		fullURL = baseURL + publicPath
	}

	return fullURL, nil
}

// deleteOldFile removes an old file from storage when a new one is uploaded.
func (s *service) deleteOldFile(fileURL string) {
	if fileURL == "" {
		return
	}

	// Strip baseURL if present
	path := fileURL
	if baseURL != "" {
		path = strings.TrimPrefix(path, baseURL)
	}

	// Convert /media/... to storage/...
	relPath := strings.TrimPrefix(path, "/media/")
	systemPath := filepath.Join(s.storage, relPath)

	if err := os.Remove(systemPath); err != nil && !os.IsNotExist(err) {
		fmt.Printf("failed to delete old file: %v\n", err)
	}
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

	// Handle avatar file upload
	if req.AvatarFile != nil {
		avatarURL, err := s.uploadFile(req.AvatarFile, "avatars", allowedImageExts, 5*1024*1024)
		if err != nil {
			return nil, errors.New("avatar upload failed: " + err.Error())
		}
		// Delete old avatar file if exists
		s.deleteOldFile(profile.AvatarURL)
		profile.AvatarURL = avatarURL
	}
	// If no new avatar file, keep existing AvatarURL unchanged

	// Handle resume file upload
	if req.ResumeFile != nil {
		resumeURL, err := s.uploadFile(req.ResumeFile, "resumes", allowedResumeExts, 10*1024*1024)
		if err != nil {
			return nil, errors.New("resume upload failed: " + err.Error())
		}
		// Delete old resume file if exists
		s.deleteOldFile(profile.ResumeURL)
		profile.ResumeURL = resumeURL
	}
	// If no new resume file, keep existing ResumeURL unchanged

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

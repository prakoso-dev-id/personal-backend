package images

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
	UploadFile(file *multipart.FileHeader) (*ImageUploadResult, error)
	DeleteImage(id uuid.UUID) error
}

type service struct {
	repo    Repository
	storage string
}

func NewService(repo Repository) Service {
	// Base storage path: relative to execution or absolute
	// Using "storage" folder in current working directory for simplicity
	return &service{
		repo:    repo,
		storage: "storage",
	}
}

func (s *service) UploadFile(file *multipart.FileHeader) (*ImageUploadResult, error) {
	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return nil, errors.New("invalid file type (only jpg, png, webp allowed)")
	}

	// Validate file size (e.g., max 5MB)
	if file.Size > 5*1024*1024 {
		return nil, errors.New("file too large (max 5MB)")
	}

	// Create directory structure: storage/uploads/
	uploadDir := filepath.Join(s.storage, "uploads")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate safe filename
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dstPath := filepath.Join(uploadDir, newFilename)

	// Save file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	// Public path for serving
	// /media/uploads/{filename}
	publicPath := fmt.Sprintf("/media/uploads/%s", newFilename)

	// Build full URL for the response
	fullURL := publicPath
	if baseURL != "" {
		fullURL = baseURL + publicPath
	}

	return &ImageUploadResult{
		FileName: newFilename,
		FilePath: fullURL,
		MimeType: file.Header.Get("Content-Type"),
		Size:     file.Size,
	}, nil
}

func (s *service) DeleteImage(id uuid.UUID) error {
	image, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if image == nil {
		return errors.New("image not found")
	}

	// Delete file
	// Public Path: /media/...
	// System Path: storage/...
	relPath := strings.TrimPrefix(image.FilePath, "/media/")
	// If it was in uploads/ (new way) or posts/uuid/ (old way), relPath handles it relative to storage root if struct match
	// The current service initializes S.storage as "storage".
	// So storage + relPath should work for both:
	// "storage" + "uploads/filename"
	// "storage" + "posts/uuid/filename"
	systemPath := filepath.Join(s.storage, relPath)

	if err := os.Remove(systemPath); err != nil && !os.IsNotExist(err) {
		// Log error but continue to delete from DB
		fmt.Printf("failed to delete file: %v\n", err)
	}

	return s.repo.Delete(id)
}

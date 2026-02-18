package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prakoso-id/personal-backend/internal/config"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(email, password string) (string, *User, error)
	UpdateEmail(userID uuid.UUID, newEmail string) error
	UpdatePassword(userID uuid.UUID, newPassword string) error
}

type service struct {
	repo Repository
	cfg  *config.Config
}

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{repo: repo, cfg: cfg}
}

func (s *service) Login(email, password string) (string, *User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *service) generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Hour * time.Duration(s.cfg.JWT.Expiration)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}

func (s *service) UpdateEmail(userID uuid.UUID, newEmail string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Check if email already taken
	existingUser, err := s.repo.FindByEmail(newEmail)
	if err != nil {
		return err
	}
	if existingUser != nil && existingUser.ID != userID {
		return errors.New("email already in use")
	}

	user.Email = newEmail
	return s.repo.Update(user)
}

func (s *service) UpdatePassword(userID uuid.UUID, newPassword string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	return s.repo.Update(user)
}

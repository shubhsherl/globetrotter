package services

import (
	"errors"

	"github.com/shubhsherl/globetrotter/backend/db"
	"github.com/shubhsherl/globetrotter/backend/models"
)

// UserService handles user-related operations
type UserService struct {
	db *db.Database
}

// NewUserService creates a new user service
func NewUserService(database *db.Database) *UserService {
	return &UserService{
		db: database,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username string) (models.User, error) {
	existingUser, err := s.db.GetUserByUsername(username)
	if err == nil && existingUser.Username != "" {
		return existingUser, errors.New("username already exists")
	}

	// Implementation depends on your database structure
	user := models.User{
		Username: username,
	}

	err = s.db.SaveUser(user)
	return user, err
}

// GetUser retrieves a user by username
func (s *UserService) GetUser(username string) (models.User, error) {
	return s.db.GetUserByUsername(username)
}

// Add other user-related methods as needed

package services

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ryanpujo/blog-app/internal/repositories"
	"github.com/ryanpujo/blog-app/models"
	"github.com/ryanpujo/blog-app/utils"
)

// UserService defines the operations available on a user service.
type UserService interface {
	Create(payload models.UserPayload) (*uint, error)
	FindById(id uint) (*models.User, error)
	FindUsers() ([]*models.User, error)
	DeleteById(id uint) error
	Update(id uint, payload *models.UserPayload) error
}

// userService implements UserService with a repository layer.
type userService struct {
	repo repositories.UserRepository
}

// NewUserService creates a new instance of userService with the given repository.
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

// Create hashes the user's password and creates a new user record.
func (s *userService) Create(payload models.UserPayload) (*uint, error) {
	isExists := s.repo.CheckIfEmailOrUsernameExist(payload.Email, payload.Username)
	if isExists {
		return nil, utils.NewDBError(utils.ErrCodeUniqueViolation, "user with a given email or username already exist", &pgconn.PgError{Code: utils.ErrCodeUniqueViolation})
	}

	hash, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}
	payload.Password = hash

	return s.repo.Create(payload)
}

// FindById retrieves a user by their ID.
func (s *userService) FindById(id uint) (*models.User, error) {
	return s.repo.FindById(id)
}

// FindUsers retrieves all users.
func (s *userService) FindUsers() ([]*models.User, error) {
	return s.repo.FindUsers()
}

// DeleteById removes a user by their ID.
func (s *userService) DeleteById(id uint) error {
	return s.repo.DeleteById(id)
}

// Update modifies an existing user record with new data.
func (s *userService) Update(id uint, payload *models.UserPayload) error {
	return s.repo.Update(id, payload)
}

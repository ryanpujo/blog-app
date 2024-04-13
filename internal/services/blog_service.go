package services

import (
	"github.com/ryanpujo/blog-app/internal/repositories"
	"github.com/ryanpujo/blog-app/models"
)

type BlogService interface {
	Create(payload models.BlogPayload) (*uint, error)
	FindById(id uint) (*models.Blog, error)
	FindBlogs() ([]*models.Blog, error)
	DeleteById(id uint) error
	Update(id uint, payload models.BlogPayload) error
}

type blogService struct {
	repo repositories.BlogRepository
}

func NewBlogService(repo repositories.BlogRepository) *blogService {
	return &blogService{
		repo: repo,
	}
}

func (s *blogService) Create(payload models.BlogPayload) (*uint, error) {
	return s.repo.Create(payload)
}

func (s *blogService) FindById(id uint) (*models.Blog, error) {
	return s.repo.FindById(id)
}

func (s *blogService) FindBlogs() ([]*models.Blog, error) {
	return s.repo.FindBlogs()
}

func (s *blogService) DeleteById(id uint) error {
	return s.repo.DeleteById(id)
}

func (s *blogService) Update(id uint, payload models.BlogPayload) error {
	return s.repo.Update(id, payload)
}

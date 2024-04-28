package services

import (
	"github.com/ryanpujo/blog-app/internal/repositories"
	"github.com/ryanpujo/blog-app/models"
)

type BlogService interface {
	Create(payload models.StoryPayload) (*uint, error)
	FindById(id uint) (*models.Story, error)
	FindBlogs() ([]*models.Story, error)
	DeleteById(id uint) error
	Update(id uint, payload models.StoryPayload) error
}

type blogService struct {
	repo repositories.StoryRepository
}

func NewBlogService(repo repositories.StoryRepository) *blogService {
	return &blogService{
		repo: repo,
	}
}

func (s *blogService) Create(payload models.StoryPayload) (*uint, error) {
	return s.repo.Create(payload)
}

func (s *blogService) FindById(id uint) (*models.Story, error) {
	return s.repo.FindById(id)
}

func (s *blogService) FindBlogs() ([]*models.Story, error) {
	return s.repo.FindBlogs()
}

func (s *blogService) DeleteById(id uint) error {
	return s.repo.DeleteById(id)
}

func (s *blogService) Update(id uint, payload models.StoryPayload) error {
	return s.repo.Update(id, payload)
}

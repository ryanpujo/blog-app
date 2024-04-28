package services

import (
	"github.com/ryanpujo/blog-app/internal/repositories"
	"github.com/ryanpujo/blog-app/models"
	"github.com/ryanpujo/blog-app/utils"
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
	payload.WordCount = utils.CountWords(payload.Content)
	if err := models.IsValidWordCountForStoryType(payload.Type, payload.WordCount); err != nil {
		return nil, err
	}
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
	payload.WordCount = utils.CountWords(payload.Content)
	if err := models.IsValidWordCountForStoryType(payload.Type, payload.WordCount); err != nil {
		return err
	}
	return s.repo.Update(id, payload)
}

package services

import (
	"github.com/ryanpujo/blog-app/internal/repositories"
	"github.com/ryanpujo/blog-app/models"
	"github.com/ryanpujo/blog-app/utils"
)

type StoryService interface {
	Create(payload models.StoryPayload) (*uint, error)
	FindById(id uint) (*models.Story, error)
	FindStories() ([]*models.Story, error)
	DeleteById(id uint) error
	Update(id uint, payload models.StoryPayload) error
}

type storyService struct {
	repo repositories.StoryRepository
}

func NewStoryService(repo repositories.StoryRepository) *storyService {
	return &storyService{
		repo: repo,
	}
}

func (s *storyService) Create(payload models.StoryPayload) (*uint, error) {
	payload.WordCount = utils.CountWords(payload.Content)
	if err := models.IsValidWordCountForStoryType(payload.Type, payload.WordCount); err != nil {
		return nil, err
	}
	return s.repo.Create(payload)
}

func (s *storyService) FindById(id uint) (*models.Story, error) {
	return s.repo.FindById(id)
}

func (s *storyService) FindStories() ([]*models.Story, error) {
	return s.repo.FindBlogs()
}

func (s *storyService) DeleteById(id uint) error {
	return s.repo.DeleteById(id)
}

func (s *storyService) Update(id uint, payload models.StoryPayload) error {
	payload.WordCount = utils.CountWords(payload.Content)
	if err := models.IsValidWordCountForStoryType(payload.Type, payload.WordCount); err != nil {
		return err
	}
	return s.repo.Update(id, payload)
}

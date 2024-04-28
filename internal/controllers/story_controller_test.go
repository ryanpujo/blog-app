package controllers_test

import (
	"github.com/ryanpujo/blog-app/models"
	"github.com/stretchr/testify/mock"
)

type MockBlogService struct {
	mock.Mock
}

func (m *MockBlogService) Create(payload models.StoryPayload) (*uint, error) {
	args := m.Called(payload)
	return args.Get(0).(*uint), args.Error(1)
}

func (m *MockBlogService) FindById(id uint) (*models.Story, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Story), args.Error(1)
}

func (m *MockBlogService) FindBlogs() ([]*models.Story, error) {
	args := m.Called()
	return args.Get(0).([]*models.Story), args.Error(1)
}

func (m *MockBlogService) DeleteById(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBlogService) Update(id uint, payload models.StoryPayload) error {
	args := m.Called(id, payload)
	return args.Error(0)
}

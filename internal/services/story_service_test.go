package services_test

import (
	"errors"
	"testing"

	"github.com/ryanpujo/blog-app/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockBlogRepository struct {
	mock.Mock
}

func (m *MockBlogRepository) Create(blog models.StoryPayload) (*uint, error) {
	args := m.Called(blog)
	return args.Get(0).(*uint), args.Error(1)
}

func (m *MockBlogRepository) FindById(id uint) (*models.Story, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Story), args.Error(1)
}

func (m *MockBlogRepository) FindBlogs() ([]*models.Story, error) {
	args := m.Called()
	return args.Get(0).([]*models.Story), args.Error(1)
}

func (m *MockBlogRepository) DeleteById(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBlogRepository) Update(id uint, payload models.StoryPayload) error {
	args := m.Called(id, payload)
	return args.Error(0)
}

var id = uint(1)

func Test_blogService_Create(t *testing.T) {
	testingTable := map[string]struct {
		payload models.StoryPayload
		arrange func()
		assert  func(t *testing.T, actualID *uint, err error)
	}{
		"success": {
			payload: models.StoryPayload{Type: 1, Content: loremGenerator.Generate(3000)},
			arrange: func() {
				mockBlogRepo.On("Create", mock.Anything).Return(&id, nil).Once()
			},
			assert: func(t *testing.T, actualID *uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), *actualID)
			},
		},
		"failed": {
			payload: models.StoryPayload{Type: 1, Content: loremGenerator.Generate(3000)},
			arrange: func() {
				mockBlogRepo.On("Create", mock.Anything).Return((*uint)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actualID *uint, err error) {
				require.Error(t, err)
				require.Nil(t, actualID)
				require.Equal(t, "failed", err.Error())
			},
		},
		"word count failed": {
			payload: models.StoryPayload{Type: 3, Content: loremGenerator.Generate(1000)},
			arrange: func() {},
			assert: func(t *testing.T, actualID *uint, err error) {
				require.Error(t, err)
				require.Nil(t, actualID)
				require.Equal(t, "story error: word count for novella should be between 20,000 and 40,000 (story type: novella, word count: 1000)", err.Error())
			},
		},
	}

	for name, tc := range testingTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			id, err := blogService.Create(tc.payload)

			tc.assert(t, id, err)
		})
	}
}

func Test_blogService_FindById(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actualBlog *models.Story, err error)
	}{
		"success": {
			arrange: func() {
				mockBlogRepo.On("FindById", mock.Anything).Return(&models.Story{Title: "test"}, nil).Once()
			},
			assert: func(t *testing.T, actualBlog *models.Story, err error) {
				require.NoError(t, err)
				require.Equal(t, "test", actualBlog.Title)
			},
		},
		"failed": {
			arrange: func() {
				mockBlogRepo.On("FindById", mock.Anything).Return((*models.Story)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actualBlog *models.Story, err error) {
				require.Error(t, err)
				require.Equal(t, "failed", err.Error())
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			blog, err := blogService.FindById(1)

			tc.assert(t, blog, err)
		})
	}
}

func Test_blogService_FindBlogs(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actualBlog []*models.Story, err error)
	}{
		"success": {
			arrange: func() {
				mockBlogRepo.On("FindBlogs").Return([]*models.Story{{Title: "test"}, {}}, nil).Once()
			},
			assert: func(t *testing.T, actualBlog []*models.Story, err error) {
				require.NoError(t, err)
				require.Equal(t, 2, len(actualBlog))
			},
		},
		"failed": {
			arrange: func() {
				mockBlogRepo.On("FindBlogs").Return(([]*models.Story)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actualBlog []*models.Story, err error) {
				require.Error(t, err)
				require.Equal(t, "failed", err.Error())
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			blogs, err := blogService.FindStories()

			tc.assert(t, blogs, err)
		})
	}
}

func Test_blogService_DeleteById(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, err error)
	}{
		"success": {
			arrange: func() {
				mockBlogRepo.On("DeleteById", mock.Anything).Return(nil).Once()
			},
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"failed": {
			arrange: func() {
				mockBlogRepo.On("DeleteById", mock.Anything).Return(errors.New("failed")).Once()
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, "failed", err.Error())
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			err := blogService.DeleteById(1)

			tc.assert(t, err)
		})
	}
}

func Test_blogService_Update(t *testing.T) {
	testTable := map[string]struct {
		payload models.StoryPayload
		arrange func()
		assert  func(t *testing.T, err error)
	}{
		"success": {
			payload: models.StoryPayload{Type: 1, Content: loremGenerator.Generate(3000)},
			arrange: func() {
				mockBlogRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"failed": {
			payload: models.StoryPayload{Type: 1, Content: loremGenerator.Generate(3000)},
			arrange: func() {
				mockBlogRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("failed")).Once()
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, "failed", err.Error())
			},
		},
		"word count failed": {
			payload: models.StoryPayload{Type: 1, Content: loremGenerator.Generate(500)},
			arrange: func() {},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, "story error: word count for short story should be between 1000 and 7500 (story type: short_story, word count: 500)", err.Error())
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			err := blogService.Update(1, tc.payload)

			tc.assert(t, err)
		})
	}
}

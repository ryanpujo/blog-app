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
		arrange func()
		assert  func(t *testing.T, actualID *uint, err error)
	}{
		"success": {
			arrange: func() {
				mockBlogRepo.On("Create", mock.Anything).Return(&id, nil).Once()
			},
			assert: func(t *testing.T, actualID *uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), *actualID)
			},
		},
		"failed": {
			arrange: func() {
				mockBlogRepo.On("Create", mock.Anything).Return((*uint)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actualID *uint, err error) {
				require.Error(t, err)
				require.Nil(t, actualID)
				require.Equal(t, "failed", err.Error())
			},
		},
	}

	for name, tc := range testingTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			id, err := blogService.Create(models.StoryPayload{})

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

			blogs, err := blogService.FindBlogs()

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
		arrange func()
		assert  func(t *testing.T, err error)
	}{
		"success": {
			arrange: func() {
				mockBlogRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"failed": {
			arrange: func() {
				mockBlogRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("failed")).Once()
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

			err := blogService.Update(1, models.StoryPayload{})

			tc.assert(t, err)
		})
	}
}

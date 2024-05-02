package controllers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/ryanpujo/blog-app/internal/response"
	"github.com/ryanpujo/blog-app/models"
	test "github.com/ryanpujo/blog-app/test/http"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func (m *MockBlogService) FindStories() ([]*models.Story, error) {
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

const storyBaseRoute = "/api/story"

var excerpt = "test excerpt"

var storyPayload = models.StoryPayload{
	Title:   "test title",
	Content: "test content",
	Slug:    "test-title",
	Excerpt: &excerpt,
	Type:    models.Novelette,
}

var storyBadPayload = models.StoryPayload{
	Title:   "test title",
	Slug:    "test-title",
	Excerpt: &excerpt,
	Type:    models.Novelette,
}

var storyTest = models.Story{
	ID:      1,
	Title:   "title test",
	Content: "content test",
}

func Test_Create_Story(t *testing.T) {
	successRet := uint(1)
	payload, _ := json.Marshal(storyPayload)
	badStoryPayload, _ := json.Marshal(storyBadPayload)
	testTbale := map[string]struct {
		uri     string
		json    []byte
		arrange func()
		assert  func(t *testing.T, statusCode int, json *response.Response)
	}{
		"success": {
			uri:  "/create/1",
			json: payload,
			arrange: func() {
				mockStoryService.On("Create", mock.Anything).Return(&successRet, nil).Once()
			},
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusCreated, statusCode)
				require.NotNil(t, json)
				require.True(t, json.Success)
				require.Equal(t, float64(1), json.Data.(map[string]any)["id"])
			},
		},
		"failed": {
			uri:  "/create/2",
			json: payload,
			arrange: func() {
				mockStoryService.On("Create", mock.Anything).Return((*uint)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.False(t, json.Success)
				require.Nil(t, json.Data)
				require.Equal(t, "An unexpected error occurred", json.Message)
			},
		},
		"validation failed": {
			uri:     "/create/1",
			json:    badStoryPayload,
			arrange: func() {},
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "The Content field is required", json.Message)
			},
		},
		"uri failed": {
			uri:     "/create/0",
			json:    payload,
			arrange: func() {},
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "The ID field must be grater than 0", json.Message)
			},
		},
	}

	for name, tc := range testTbale {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			res, code, err := test.NewHttpTest(http.MethodPost, tc.uri, test.WithBaseUri(storyBaseRoute), test.WithJson(tc.json)).
				ExecuteTest(mux)
			require.NoError(t, err)

			tc.assert(t, code, res)
		})
	}
}
func Test_Find_Story(t *testing.T) {
	testTable := map[string]struct {
		uri     string
		arrange func()
		assert  func(t *testing.T, statusCode int, json *response.Response)
	}{
		"success": {
			uri: "/1",
			arrange: func() {
				mockStoryService.On("FindById", mock.Anything).Return(&storyTest, nil).Once()
			},
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.NotNil(t, json.Data)
				require.Equal(t, storyTest.Title, json.Data.(map[string]any)["story"].(map[string]any)["title"])
				mockStoryService.AssertCalled(t, "FindById", uint(1))
			},
		},
		"failed": {
			uri: "/1",
			arrange: func() {
				mockStoryService.On("FindById", mock.Anything).Return((*models.Story)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "An unexpected error occurred", json.Message)
			},
		},
		"validation failed": {
			uri:     "/0",
			arrange: func() {},
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "The StoryID field must be grater than 0", json.Message)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			res, code, err := test.NewHttpTest(http.MethodGet, tc.uri, test.WithBaseUri(storyBaseRoute)).ExecuteTest(mux)
			require.NoError(t, err)

			tc.assert(t, code, res)
		})
	}
}

func Test_Find_Stories(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, statusCode int, res *response.Response)
	}{
		"success": {
			arrange: func() {
				mockStoryService.On("FindStories", mock.Anything).Return([]*models.Story{{}, {}, {}}, nil).Once()
			},
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.NotNil(t, res.Data)
				require.Equal(t, 3, len(res.Data.(map[string]any)["stories"].([]any)))
			},
		},
		"failed": {
			arrange: func() {
				mockStoryService.On("FindStories", mock.Anything).Return(([]*models.Story)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, res.Data)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			res, code, err := test.NewHttpTest(http.MethodGet, "/", test.WithBaseUri(storyBaseRoute)).ExecuteTest(mux)
			require.NoError(t, err)

			tc.assert(t, code, res)
		})
	}
}
func Test_Update_Story(t *testing.T) {
	payload, _ := json.Marshal(storyPayload)
	badStoryPayload, _ := json.Marshal(storyBadPayload)
	testTable := map[string]struct {
		uri     string
		json    []byte
		arrange func()
		assert  func(t *testing.T, statusCode int, res *response.Response)
	}{
		"success": {
			uri:  "/1/user/1",
			json: payload,
			arrange: func() {
				mockStoryService.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.Nil(t, res)
			},
		},
		"failed": {
			uri:  "/1/user/1",
			json: payload,
			arrange: func() {
				mockStoryService.On("Update", mock.Anything, mock.Anything).Return(errors.New("failed")).Once()
			},
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, res)
				require.Equal(t, "An unexpected error occurred", res.Message)
			},
		},
		"story uri failed": {
			uri:  "/0/user/1",
			json: payload,
			arrange: func() {
			},
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, res)
				require.Equal(t, "The StoryID field must be grater than 0", res.Message)
			},
		},
		"user uri failed": {
			uri:  "/1/user/0",
			json: payload,
			arrange: func() {
			},
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, res)
				require.Equal(t, "The ID field must be grater than 0", res.Message)
			},
		},
		"json failed": {
			uri:  "/1/user/1",
			json: badStoryPayload,
			arrange: func() {
			},
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, res)
				require.Equal(t, "The Content field is required", res.Message)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			res, code, err := test.NewHttpTest(http.MethodPatch, tc.uri, test.WithBaseUri(storyBaseRoute), test.WithJson(tc.json)).ExecuteTest(mux)
			require.NoError(t, err)

			tc.assert(t, code, res)
		})
	}
}

func Test_Delete_Story(t *testing.T) {
	testTable := map[string]struct {
		uri     string
		arrange func()
		assert  func(t *testing.T, statusCode int, res *response.Response)
	}{
		"success": {
			uri: "/1",
			arrange: func() {
				mockStoryService.On("DeleteById", mock.Anything).Return(nil).Once()
			},
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.Nil(t, res)
			},
		},
		"failed": {
			uri: "/1",
			arrange: func() {
				mockStoryService.On("DeleteById", mock.Anything).Return(errors.New("failed")).Once()
			},
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, res)
				require.Equal(t, "An unexpected error occurred", res.Message)
			},
		},
		"uri failed": {
			uri: "/0",
			arrange: func() {
				mockStoryService.On("DeleteById", mock.Anything).Return(errors.New("failed")).Once()
			},
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, res)
				require.Equal(t, "The StoryID field must be grater than 0", res.Message)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			res, code, err := test.NewHttpTest(http.MethodDelete, tc.uri, test.WithBaseUri(storyBaseRoute)).ExecuteTest(mux)
			require.NoError(t, err)

			tc.assert(t, code, res)
		})
	}
}

package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ryanpujo/blog-app/internal/response"
	"github.com/ryanpujo/blog-app/models"
	test "github.com/ryanpujo/blog-app/test/http"
	"github.com/stretchr/testify/require"
)

const storyBaseRoute = "/api/story"

var excerpt = "test excerpt"

var storyPayload = models.StoryPayload{
	Title:   "test title",
	Slug:    "test-title",
	Excerpt: &excerpt,
	Type:    models.Novelette,
}

var uniquenessViolatonPayload = models.StoryPayload{
	Title:   "test title",
	Slug:    "tenth-blog-post",
	Excerpt: &excerpt,
	Type:    models.FlashFiction,
}

var storyContentFailedPayload = models.StoryPayload{
	Title:   "test title",
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

func Test_Create_Story(t *testing.T) {
	storyPayload.Content = loremGenerator.Generate(8000)
	payload, _ := json.Marshal(storyPayload)
	uniquenessViolatonPayload.Content = loremGenerator.Generate(200)
	uniquenessViolation, _ := json.Marshal(uniquenessViolatonPayload)
	badStoryPayload, _ := json.Marshal(storyBadPayload)
	storyContentFailedPayload.Content = loremGenerator.Generate(1000)
	contentFailedPayload, _ := json.Marshal(storyContentFailedPayload)
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

			},
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusCreated, statusCode)
				require.NotNil(t, json)
				require.True(t, json.Success)
				require.Equal(t, float64(2), json.Data.(map[string]any)["id"])
			},
		},
		"uniqueness violated": {
			uri:  "/create/2",
			json: uniquenessViolation,
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.False(t, json.Success)
				require.Nil(t, json.Data)
				require.Equal(t, "user with a given email or username already exist", json.Message)
			},
		},
		"validation failed": {
			uri:  "/create/1",
			json: badStoryPayload,
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "The Content field is required", json.Message)
			},
		},
		"uri failed": {
			uri:  "/create/0",
			json: payload,
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "The ID field must be grater than 0", json.Message)
			},
		},
		"word count failed": {
			uri:  "/create/1",
			json: contentFailedPayload,
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.False(t, json.Success)
				require.Nil(t, json.Data)
				require.Equal(t, "word count for novelette should be between 7500 and 20,000", json.Message)
			},
		},
	}

	for _, tc := range testTbale {
		res, code, err := test.NewHttpTest(http.MethodPost, tc.uri, test.WithBaseUri(storyBaseRoute), test.WithJson(tc.json)).
			ExecuteTest(mux)
		require.NoError(t, err)

		tc.assert(t, code, res)
	}
}
func Test_Find_Story(t *testing.T) {
	testTable := map[string]struct {
		uri    string
		assert func(t *testing.T, statusCode int, json *response.Response)
	}{
		"success": {
			uri: "/2",
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.NotNil(t, json.Data)
				require.Equal(t, storyPayload.Title, json.Data.(map[string]any)["story"].(map[string]any)["title"])
			},
		},
		"not found": {
			uri: "/20",
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusNotFound, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "data not found", json.Message)
			},
		},
		"validation failed": {
			uri: "/0",
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "The StoryID field must be grater than 0", json.Message)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
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
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.NotNil(t, res.Data)
				require.Equal(t, 2, len(res.Data.(map[string]any)["stories"].([]any)))
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			res, code, err := test.NewHttpTest(http.MethodGet, "/", test.WithBaseUri(storyBaseRoute)).ExecuteTest(mux)
			require.NoError(t, err)

			tc.assert(t, code, res)
		})
	}
}

func Test_Update_Story(t *testing.T) {
	payload, _ := json.Marshal(storyPayload)
	badStoryPayload, _ := json.Marshal(storyBadPayload)
	storyContentFailedPayload.Content = loremGenerator.Generate(1000)
	contentFailedPayload, _ := json.Marshal(storyContentFailedPayload)
	testTable := map[string]struct {
		uri    string
		json   []byte
		assert func(t *testing.T, statusCode int, res *response.Response)
	}{
		"success": {
			uri:  "/2/user/1",
			json: payload,
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.Nil(t, res)
			},
		},
		"word count failed": {
			uri:  "/2/user/1",
			json: contentFailedPayload,
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, res)
				require.Equal(t, "word count for novelette should be between 7500 and 20,000", res.Message)
			},
		},
		"story uri failed": {
			uri:  "/0/user/1",
			json: payload,
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, res)
				require.Equal(t, "The StoryID field must be grater than 0", res.Message)
			},
		},
		"user uri failed": {
			uri:  "/1/user/0",
			json: payload,
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, res)
				require.Equal(t, "The ID field must be grater than 0", res.Message)
			},
		},
		"json failed": {
			uri:  "/1/user/1",
			json: badStoryPayload,
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, res)
				require.Equal(t, "The Content field is required", res.Message)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			res, code, err := test.NewHttpTest(http.MethodPatch, tc.uri, test.WithBaseUri(storyBaseRoute), test.WithJson(tc.json)).ExecuteTest(mux)
			require.NoError(t, err)

			tc.assert(t, code, res)
		})
	}
}

func Test_Delete_Story(t *testing.T) {
	testTable := map[string]struct {
		uri    string
		assert func(t *testing.T, statusCode int, res *response.Response)
	}{
		"success": {
			uri: "/1",
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.Nil(t, res)
			},
		},
		"data not found": {
			uri: "/20",
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusNotFound, statusCode)
				require.NotNil(t, res)
				require.Equal(t, "data not found", res.Message)
			},
		},
		"uri failed": {
			uri: "/0",
			assert: func(t *testing.T, statusCode int, res *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, res)
				require.Equal(t, "The StoryID field must be grater than 0", res.Message)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			res, code, err := test.NewHttpTest(http.MethodDelete, tc.uri, test.WithBaseUri(storyBaseRoute)).ExecuteTest(mux)
			require.NoError(t, err)

			tc.assert(t, code, res)
		})
	}
}

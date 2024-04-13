package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ryanpujo/blog-app/internal/response"
	"github.com/ryanpujo/blog-app/models"
	test "github.com/ryanpujo/blog-app/test/http"
	"github.com/stretchr/testify/require"
)

var (
	payload = models.UserPayload{
		FirstName: "michael",
		LastName:  "townley",
		Username:  "desanta",
		Password:  "fucktrevor",
		Email:     "townley@gmail.com",
	}
	badPayload = models.UserPayload{
		FirstName: "michael",
		LastName:  "townley",
		Username:  "desanta",
		Password:  "fucktrevor",
		Email:     "townleygmail.com",
	}
)

var userBaseUri = "/api/user"

func Test_Create(t *testing.T) {
	jsonPayload, _ := json.Marshal(payload)
	badJson, _ := json.Marshal(badPayload)
	testTable := map[string]struct {
		Json    []byte
		Arrange func()
		Assert  func(t *testing.T, statusCode int, json *response.Response)
	}{
		"success": {
			Json: jsonPayload,
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusCreated, statusCode)
				require.NotNil(t, json)
				require.Equal(t, float64(11), json.Data.(map[string]any)["id"])
			},
		},
		"failed to create": {
			Json: jsonPayload,
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Equal(t, "user with a given email or username already exist", json.Message)
			},
		},
		"validation error": {
			Json: badJson,
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Equal(t, "The Email field must be a valid email address", json.Message)
			},
		},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			jsonRes, code := test.NewHttpTest(http.MethodPost, "/create", test.WithBaseUri(userBaseUri), test.WithJson(testCase.Json)).
				ExecuteTest(t, mux)
			testCase.Assert(t, code, jsonRes)
		})
	}
}

func Test_FindById(t *testing.T) {
	testTable := map[string]struct {
		ID      uint
		Arrange func()
		Assert  func(t *testing.T, statusCode int, json *response.Response)
	}{
		"success": {
			ID: 1,
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.NotNil(t, json.Data)
				require.Equal(t, "johndoe", json.Data.(map[string]any)["user"].(map[string]any)["username"])
			},
		},
		"0 id": {
			ID: 0,
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "The ID field must be grater than 0", json.Message)
			},
		},
		"failed": {
			ID: 40,
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusNotFound, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "data not found", json.Message)
			},
		},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			res, code := test.NewHttpTest(http.MethodGet, fmt.Sprintf("/%d", testCase.ID), test.WithBaseUri(userBaseUri)).
				ExecuteTest(t, mux)

			testCase.Assert(t, code, res)
		})
	}
}

func Test_FindUsers(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, statusCode int, json *response.Response)
	}{
		"success": {
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.NotNil(t, json.Data)
				require.Equal(t, 11, len(json.Data.(map[string]any)["users"].([]any)))
			},
		},
	}
	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			res, code := test.NewHttpTest(http.MethodGet, "/", test.WithBaseUri(userBaseUri)).ExecuteTest(t, mux)

			testCase.assert(t, code, res)
		})
	}
}

func Test_DeleteById(t *testing.T) {
	testTable := map[string]struct {
		ID      uint
		Arrange func()
		Assert  func(t *testing.T, statusCode int, json *response.Response)
	}{
		"success": {
			ID: 8,
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.Nil(t, json)
			},
		},
		"0 id": {
			ID: 0,
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "The ID field must be grater than 0", json.Message)
			},
		},
		"failed": {
			ID: 8,
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusNotFound, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "data not found", json.Message)
			},
		},
	}

	for _, testCase := range testTable {
		res, code := test.NewHttpTest(http.MethodDelete, fmt.Sprintf("/%d", testCase.ID), test.WithBaseUri(userBaseUri)).
			ExecuteTest(t, mux)

		testCase.Assert(t, code, res)
	}
}

func Test_Update(t *testing.T) {
	var updated = models.UserPayload{
		FirstName: "john",
		LastName:  "doe",
		Username:  "janedoe",
		Password:  "password123",
		Email:     "john.doe@example.com",
	}
	jsonPayload, _ := json.Marshal(updated)
	badJson, _ := json.Marshal(badPayload)
	testTable := map[string]struct {
		ID      uint
		JSON    []byte
		arrange func()
		assert  func(t *testing.T, statusCode int, json *response.Response)
	}{
		"success": {
			ID:   1,
			JSON: jsonPayload,
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.Nil(t, json)
			},
		},
		"failed": {
			ID:   32,
			JSON: jsonPayload,
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusNotFound, statusCode)
				require.NotNil(t, json)
				require.Nil(t, json.Data)
				require.Equal(t, "data not found", json.Message)
			},
		},
		"bad json": {
			ID:   1,
			JSON: badJson,
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, json)
				require.Nil(t, json.Data)
				require.Equal(t, "The Email field must be a valid email address", json.Message)
			},
		},
		"bad uri": {
			ID:      0,
			JSON:    jsonPayload,
			arrange: func() {},
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, json)
				require.Nil(t, json.Data)
				require.Equal(t, "The ID field must be grater than 0", json.Message)
			},
		},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			res, code := test.NewHttpTest(http.MethodPatch, fmt.Sprintf("/%d", testCase.ID), test.WithBaseUri(userBaseUri), test.WithJson(testCase.JSON)).
				ExecuteTest(t, mux)

			testCase.assert(t, code, res)
		})
	}
}

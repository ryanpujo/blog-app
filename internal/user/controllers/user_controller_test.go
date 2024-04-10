package controllers_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/blog-app/internal/adapter"
	"github.com/ryanpujo/blog-app/internal/response"
	"github.com/ryanpujo/blog-app/internal/route"
	"github.com/ryanpujo/blog-app/internal/user/controllers"
	"github.com/ryanpujo/blog-app/models"
	test "github.com/ryanpujo/blog-app/test/http"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(payload models.UserPayload) (uint, error) {
	args := m.Called(payload)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockUserService) FindById(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) FindUsers() ([]*models.User, error) {
	args := m.Called()
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserService) DeleteById(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) Update(id uint, payload *models.UserPayload) error {
	args := m.Called(id, payload)
	return args.Error(0)
}

var (
	mockService *MockUserService
	mux         *gin.Engine
	baseUri     = "/api/user"
	payload     = models.UserPayload{
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

func TestMain(m *testing.M) {
	mockService = new(MockUserService)
	userController := controllers.NewUserController(mockService)
	mux = route.Route(adapter.AppController{UserController: userController})
	os.Exit(m.Run())
}

func Test_userController_Create(t *testing.T) {
	jsonPayload, _ := json.Marshal(payload)
	badJson, _ := json.Marshal(badPayload)
	testTable := map[string]struct {
		Json    []byte
		Arrange func()
		Assert  func(t *testing.T, statusCode int, json *response.Response)
	}{
		"success": {
			Json: jsonPayload,
			Arrange: func() {
				mockService.On("Create", mock.Anything).Return(uint(1), nil).Once()
			},
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusCreated, statusCode)
				require.NotNil(t, json)
				require.Equal(t, float64(1), json.Data.(map[string]any)["id"])
			},
		},
		"failed": {
			Json: jsonPayload,
			Arrange: func() {
				mockService.On("Create", mock.Anything).Return(uint(0), errors.New("failed")).Once()
			},
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Equal(t, "An unexpected error occurred", json.Message)
			},
		},
		"validation error": {
			Json: badJson,
			Arrange: func() {
				// mockService.On("Create", mock.Anything).Return(uint(0), errors.New("failed")).Once()
			},
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Equal(t, "The Email field must be a valid email address", json.Message)
			},
		},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			testCase.Arrange()
			jsonRes, code := test.NewHttpTest(http.MethodPost, "/create", test.WithBaseUri(baseUri), test.WithJson(testCase.Json)).
				ExecuteTest(t, mux)
			testCase.Assert(t, code, jsonRes)
		})
	}
}

func Test_userController_FindById(t *testing.T) {
	testTable := map[string]struct {
		ID      uint
		Arrange func()
		Assert  func(t *testing.T, statusCode int, json *response.Response)
	}{
		"success": {
			ID: 1,
			Arrange: func() {
				mockService.On("FindById", mock.Anything).Return(&models.User{Username: "okeoke"}, nil).Once()
			},
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.NotNil(t, json.Data)
				require.Equal(t, "okeoke", json.Data.(map[string]any)["user"].(map[string]any)["username"])
				mockService.AssertCalled(t, "FindById", uint(1))
			},
		},
		"0 id": {
			ID: 0,
			Arrange: func() {
			},
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "The ID field must be grater than 0", json.Message)
			},
		},
		"failed": {
			ID: 1,
			Arrange: func() {
				mockService.On("FindById", mock.Anything).Return((*models.User)(nil), errors.New("failed")).Once()
			},
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "An unexpected error occurred", json.Message)
			},
		},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			testCase.Arrange()

			res, code := test.NewHttpTest(http.MethodGet, fmt.Sprintf("/%d", testCase.ID), test.WithBaseUri(baseUri)).
				ExecuteTest(t, mux)

			testCase.Assert(t, code, res)
		})
	}
}

func Test_userController_FindUsers(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, statusCode int, json *response.Response)
	}{
		"success": {
			arrange: func() {
				mockService.On("FindUsers").Return([]*models.User{{}, {}}, nil).Once()
			},
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.NotNil(t, json.Data)
				require.Equal(t, 2, len(json.Data.(map[string]any)["users"].([]any)))
			},
		},
		"failed": {
			arrange: func() {
				mockService.On("FindUsers").Return(([]*models.User)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "An unexpected error occurred", json.Message)
			},
		},
	}
	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			testCase.arrange()

			res, code := test.NewHttpTest(http.MethodGet, "/", test.WithBaseUri(baseUri)).ExecuteTest(t, mux)

			testCase.assert(t, code, res)
		})
	}
}

func Test_userController_DeleteById(t *testing.T) {
	testTable := map[string]struct {
		ID      uint
		Arrange func()
		Assert  func(t *testing.T, statusCode int, json *response.Response)
	}{
		"success": {
			ID: 1,
			Arrange: func() {
				mockService.On("DeleteById", mock.Anything).Return(nil).Once()
			},
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.Nil(t, json)
				mockService.AssertCalled(t, "DeleteById", uint(1))
			},
		},
		"0 id": {
			ID: 0,
			Arrange: func() {
			},
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "The ID field must be grater than 0", json.Message)
			},
		},
		"failed": {
			ID: 1,
			Arrange: func() {
				mockService.On("DeleteById", mock.Anything).Return(errors.New("failed")).Once()
			},
			Assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Nil(t, json.Data)
				require.Equal(t, "An unexpected error occurred", json.Message)
			},
		},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			testCase.Arrange()

			res, code := test.NewHttpTest(http.MethodDelete, fmt.Sprintf("/%d", testCase.ID), test.WithBaseUri(baseUri)).
				ExecuteTest(t, mux)

			testCase.Assert(t, code, res)
		})
	}
}

func Test_userController_Update(t *testing.T) {
	jsonPayload, _ := json.Marshal(payload)
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
			arrange: func() {
				mockService.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusOK, statusCode)
				require.Nil(t, json)
			},
		},
		"failed": {
			ID:   1,
			JSON: jsonPayload,
			arrange: func() {
				mockService.On("Update", mock.Anything, mock.Anything).Return(errors.New("failed")).Once()
			},
			assert: func(t *testing.T, statusCode int, json *response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotNil(t, json)
				require.Nil(t, json.Data)
				require.Equal(t, "An unexpected error occurred", json.Message)
			},
		},
		"bad json": {
			ID:      1,
			JSON:    badJson,
			arrange: func() {},
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
			testCase.arrange()

			res, code := test.NewHttpTest(http.MethodPatch, fmt.Sprintf("/%d", testCase.ID), test.WithBaseUri(baseUri), test.WithJson(testCase.JSON)).
				ExecuteTest(t, mux)

			testCase.assert(t, code, res)
		})
	}
}

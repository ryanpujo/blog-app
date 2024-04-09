package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/blog-app/internal/adapter"
	"github.com/ryanpujo/blog-app/internal/response"
	"github.com/ryanpujo/blog-app/internal/route"
	"github.com/ryanpujo/blog-app/internal/user/controllers"
	"github.com/ryanpujo/blog-app/models"
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

func (m *MockUserService) Update(payload *models.UserPayload) error {
	args := m.Called(payload)
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
		json    []byte
		arrange func()
		assert  func(t *testing.T, statusCode int, json response.Response)
	}{
		"success": {
			json: jsonPayload,
			arrange: func() {
				mockService.On("Create", mock.Anything).Return(uint(1), nil).Once()
			},
			assert: func(t *testing.T, statusCode int, json response.Response) {
				require.Equal(t, http.StatusCreated, statusCode)
				require.NotNil(t, json)
				require.Equal(t, float64(1), json.Data.(map[string]any)["id"])
			},
		},
		"failed": {
			json: jsonPayload,
			arrange: func() {
				mockService.On("Create", mock.Anything).Return(uint(0), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, statusCode int, json response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Equal(t, "An unexpected error occurred", json.Message)
			},
		},
		"validation error": {
			json: badJson,
			arrange: func() {
				// mockService.On("Create", mock.Anything).Return(uint(0), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, statusCode int, json response.Response) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.Equal(t, "The Email field must be a valid email address", json.Message)
			},
		},
	}

	for name, test := range testTable {
		t.Run(name, func(t *testing.T) {
			test.arrange()

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/create", baseUri), bytes.NewReader(test.json))
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			var jsonRes response.Response
			json.NewDecoder(rr.Body).Decode(&jsonRes)

			test.assert(t, rr.Code, jsonRes)
		})
	}
}

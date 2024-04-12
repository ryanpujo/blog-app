package services_test

import (
	"errors"
	"os"
	"testing"

	"github.com/ryanpujo/blog-app/internal/services"
	"github.com/ryanpujo/blog-app/models"
	"github.com/ryanpujo/blog-app/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// UserRepository is a mock type for the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

// Create is a mock method that simulates the Create method of the UserRepository interface
func (_m *MockUserRepository) Create(payload models.UserPayload) (*uint, error) {
	ret := _m.Called(payload)
	return ret.Get(0).(*uint), ret.Error(1)
}

// FindById is a mock method that simulates the FindById method of the UserRepository interface
func (_m *MockUserRepository) FindById(id uint) (*models.User, error) {
	ret := _m.Called(id)
	return ret.Get(0).(*models.User), ret.Error(1)
}

// FindUsers is a mock method that simulates the FindUsers method of the UserRepository interface
func (_m *MockUserRepository) FindUsers() ([]*models.User, error) {
	ret := _m.Called()
	return ret.Get(0).([]*models.User), ret.Error(1)
}

// DeleteById is a mock method that simulates the DeleteById method of the UserRepository interface
func (_m *MockUserRepository) DeleteById(id uint) error {
	ret := _m.Called(id)
	return ret.Error(0)
}

// Update is a mock method that simulates the Update method of the UserRepository interface
func (_m *MockUserRepository) Update(id uint, user *models.UserPayload) error {
	ret := _m.Called(id, user)
	return ret.Error(0)
}
func (_m *MockUserRepository) CheckIfEmailOrUsernameExist(email, username string) bool {
	ret := _m.Called(email, username)
	return ret.Bool(0)
}

// Mock for UserRepository interface
var mockRepo *MockUserRepository

// UserService instance with injected mock
var userService services.UserService

// TestMain sets up the mock repository and userService before running the tests
func TestMain(m *testing.M) {
	mockRepo = new(MockUserRepository)
	userService = services.NewUserService(mockRepo)
	os.Exit(m.Run())
}

// Test_userService_Create tests the Create method of userService
func Test_userService_Create(t *testing.T) {
	succesRet := uint(1)
	// Define a test table to run subtests
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actualID *uint, err error)
	}{
		// Subtest for successful user creation
		"success": {
			arrange: func() {
				// Arrange for a successful creation by setting up the mock expectations
				mockRepo.On("Create", mock.Anything).Return(&succesRet, nil).Once()
				utils.HashPassword = utils.EncryptPassword
				mockRepo.On("CheckIfEmailOrUsernameExist", mock.Anything, mock.Anything).Return(false).Once()
			},
			assert: func(t *testing.T, actualID *uint, err error) {
				// Assert that no error occurred and the returned ID is as expected
				require.NoError(t, err)
				require.Equal(t, &succesRet, actualID)
				mockRepo.AssertCalled(t, "Create", mock.AnythingOfType("models.UserPayload"))
			},
		},
		// Subtest for hashing error during user creation
		"hashing error": {
			arrange: func() {
				mockRepo.On("CheckIfEmailOrUsernameExist", mock.Anything, mock.Anything).Return(false).Once()
				// Arrange for a hashing error by overriding the HashPassword function
				utils.HashPassword = func(plain string) (string, error) {
					return "", errors.New("hash password")
				}
			},
			assert: func(t *testing.T, actualID *uint, err error) {
				// Assert that an error occurred and the error message is as expected
				require.Error(t, err)
				require.Equal(t, "hash password", err.Error())
				require.Zero(t, actualID)
				mockRepo.AssertNotCalled(t, "Create")
			},
		},
		// Subtest for failed user creation
		"failed": {
			arrange: func() {
				// Arrange for a failed creation by setting up the mock expectations
				mockRepo.On("Create", mock.Anything).Return((*uint)(nil), errors.New("failed to create")).Once()
				utils.HashPassword = utils.EncryptPassword
				mockRepo.On("CheckIfEmailOrUsernameExist", mock.Anything, mock.Anything).Return(false).Once()
			},
			assert: func(t *testing.T, actualID *uint, err error) {
				// Assert that an error occurred and the returned ID is zero
				require.Error(t, err)
				require.Equal(t, "failed to create", err.Error())
				require.Zero(t, actualID)
				mockRepo.AssertCalled(t, "Create", mock.AnythingOfType("models.UserPayload"))
			},
		},
		"email or username exists": {
			arrange: func() {
				mockRepo.On("CheckIfEmailOrUsernameExist", mock.Anything, mock.Anything).Return(true).Once()
			},
			assert: func(t *testing.T, actualID *uint, err error) {
				require.Zero(t, actualID)
				require.NotNil(t, err)
				require.IsType(t, utils.DBError{}, err)
			},
		},
	}

	// Run the subtests defined in the test table
	for name, test := range testTable {
		t.Run(name, func(t *testing.T) {
			test.arrange() // Arrange the test scenario

			// Act by calling the Create method
			id, err := userService.Create(models.UserPayload{})

			test.assert(t, id, err) // Assert the expected outcome
		})
	}
}

func Test_userService_FindById(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actual *models.User, err error)
	}{
		"success": {
			arrange: func() {
				mockRepo.On("FindById", mock.Anything).Return(&models.User{Username: "townley"}, nil).Once()
			},
			assert: func(t *testing.T, actual *models.User, err error) {
				require.NoError(t, err)
				require.NotNil(t, actual)
				require.Equal(t, "townley", actual.Username)
				mockRepo.AssertCalled(t, "FindById", mock.AnythingOfType("uint"))
			},
		},
		"failed": {
			arrange: func() {
				mockRepo.On("FindById", mock.Anything).Return((*models.User)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actual *models.User, err error) {
				require.Error(t, err)
				require.Nil(t, actual)
				require.Equal(t, "failed", err.Error())
				mockRepo.AssertCalled(t, "FindById", mock.AnythingOfType("uint"))
			},
		},
	}

	for name, test := range testTable {
		t.Run(name, func(t *testing.T) {
			test.arrange()

			user, err := userService.FindById(1)

			test.assert(t, user, err)
		})
	}
}

func Test_userService_FindUsers(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actual []*models.User, err error)
	}{
		"success": {
			arrange: func() {
				mockRepo.On("FindUsers").Return([]*models.User{{Username: "townley"}, {Username: "trevor"}}, nil).Once()
			},
			assert: func(t *testing.T, actual []*models.User, err error) {
				require.NoError(t, err)
				require.NotNil(t, actual)
				require.Equal(t, 2, len(actual))
				mockRepo.AssertCalled(t, "FindUsers")
			},
		},
		"failed": {
			arrange: func() {
				mockRepo.On("FindUsers").Return(([]*models.User)(nil), errors.New("failed")).Once()
			},
			assert: func(t *testing.T, actual []*models.User, err error) {
				require.Error(t, err)
				require.Nil(t, actual)
				require.Equal(t, "failed", err.Error())
				mockRepo.AssertCalled(t, "FindUsers")
			},
		},
	}

	for name, test := range testTable {
		t.Run(name, func(t *testing.T) {
			test.arrange()

			user, err := userService.FindUsers()

			test.assert(t, user, err)
		})
	}
}

func Test_userService_DeleteById(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, err error)
	}{
		"success": {
			arrange: func() {
				mockRepo.On("DeleteById", mock.Anything).Return(nil).Once()
			},
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"failed": {
			arrange: func() {
				mockRepo.On("DeleteById", mock.Anything).Return(errors.New("failed")).Once()
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, "failed", err.Error())
			},
		},
	}

	for name, test := range testTable {
		t.Run(name, func(t *testing.T) {
			test.arrange()

			err := userService.DeleteById(1)

			test.assert(t, err)
		})
	}
}

func Test_userService_Update(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, err error)
	}{
		"success": {
			arrange: func() {
				mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"failed": {
			arrange: func() {
				mockRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("failed")).Once()
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, "failed", err.Error())
			},
		},
	}

	for name, test := range testTable {
		t.Run(name, func(t *testing.T) {
			test.arrange()

			err := userService.Update(1, &models.UserPayload{})

			test.assert(t, err)
		})
	}
}

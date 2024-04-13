package services_test

import (
	"os"
	"testing"

	"github.com/ryanpujo/blog-app/internal/services"
)

var (
	mockBlogRepo *MockBlogRepository
	blogService  services.BlogService
	mockRepo     *MockUserRepository
	userService  services.UserService
)

// TestMain sets up the mock repository and userService before running the tests
func TestMain(m *testing.M) {
	mockRepo = new(MockUserRepository)
	userService = services.NewUserService(mockRepo)

	mockBlogRepo = new(MockBlogRepository)
	blogService = services.NewBlogService(mockBlogRepo)
	os.Exit(m.Run())
}

package services_test

import (
	"os"
	"testing"

	"github.com/ryanpujo/blog-app/internal/services"

	lorem "github.com/derektata/lorem/ipsum"
)

var (
	mockBlogRepo   *MockBlogRepository
	blogService    services.StoryService
	mockRepo       *MockUserRepository
	userService    services.UserService
	loremGenerator lorem.Generator
)

// TestMain sets up the mock repository and userService before running the tests
func TestMain(m *testing.M) {
	mockRepo = new(MockUserRepository)
	userService = services.NewUserService(mockRepo)

	mockBlogRepo = new(MockBlogRepository)
	blogService = services.NewStoryService(mockBlogRepo)
	loremGenerator = *lorem.NewGenerator()
	os.Exit(m.Run())
}

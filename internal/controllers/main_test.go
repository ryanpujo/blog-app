package controllers_test

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/blog-app/internal/adapter"
	"github.com/ryanpujo/blog-app/internal/controllers"
	"github.com/ryanpujo/blog-app/internal/route"
)

var (
	mockService      *MockUserService
	mockStoryService *MockBlogService
	mux              *gin.Engine
)

func TestMain(m *testing.M) {
	mockService = new(MockUserService)
	userController := controllers.NewUserController(mockService)

	mockStoryService = new(MockBlogService)
	storyController := controllers.NewStoryController(mockStoryService)

	adapter := adapter.AppController{
		UserController:  userController,
		StoryController: storyController,
	}
	mux = route.Route(adapter)
	os.Exit(m.Run())
}

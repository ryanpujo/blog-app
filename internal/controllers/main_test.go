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
	mockService     *MockUserService
	mockBlogService *MockBlogService
	mux             *gin.Engine
)

func TestMain(m *testing.M) {
	mockService = new(MockUserService)
	userController := controllers.NewUserController(mockService)

	mockBlogService = new(MockBlogService)
	blogContorller := controllers.NewBlogController(mockBlogService)

	adapter := adapter.AppController{
		UserController: userController,
		BlogController: blogContorller,
	}
	mux = route.Route(adapter)
	os.Exit(m.Run())
}

package adapter

import "github.com/ryanpujo/blog-app/internal/controllers"

type AppController struct {
	UserController controllers.UserController
	BlogController controllers.StoryController
}

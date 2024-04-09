package adapter

import "github.com/ryanpujo/blog-app/internal/user/controllers"

type AppController struct {
	UserController controllers.UserController
}

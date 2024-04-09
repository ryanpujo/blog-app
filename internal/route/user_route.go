package route

import "github.com/ryanpujo/blog-app/internal/user/controllers"

func UserRoute(uc controllers.UserController) {
	userRoute := mux.Group("/api/user")

	userRoute.POST("/create", uc.Create)
}

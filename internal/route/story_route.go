package route

import "github.com/ryanpujo/blog-app/internal/controllers"

func StoryRoute(storyController controllers.StoryController) {
	baseRoute := mux.Group("/api/story")

	baseRoute.POST("/create/:id", storyController.Create)
	baseRoute.GET("/:storyID", storyController.FindById)
	baseRoute.GET("/", storyController.FindStories)
	baseRoute.PATCH("/:storyID/user/:id", storyController.Update)
	baseRoute.DELETE("/:storyID", storyController.DeleteById)
}

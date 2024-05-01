package route

import "github.com/ryanpujo/blog-app/internal/controllers"

func StoryRoute(storyController controllers.StoryController) {
	baseRoute := mux.Group("/api/story")

	baseRoute.POST("/create", storyController.Create)
	baseRoute.GET("/:id", storyController.FindById)
	baseRoute.GET("/", storyController.FindStories)
	baseRoute.PATCH("/:id", storyController.Update)
	baseRoute.DELETE("/:id", storyController.DeleteById)
}

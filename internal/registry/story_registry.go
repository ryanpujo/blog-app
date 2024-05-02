package registry

import (
	"github.com/ryanpujo/blog-app/internal/controllers"
	"github.com/ryanpujo/blog-app/internal/repositories"
	"github.com/ryanpujo/blog-app/internal/services"
)

func (r registry) NewStoryRepository() repositories.StoryRepository {
	return repositories.NewStoryRepository(r.DB)
}

func (r registry) NewStoryService() services.StoryService {
	return services.NewStoryService(r.NewStoryRepository())
}

func (r registry) NewStoryController() controllers.StoryController {
	return controllers.NewStoryController(r.NewStoryService())
}

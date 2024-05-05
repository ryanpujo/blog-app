package registry

import (
	"github.com/ryanpujo/blog-app/database"
	"github.com/ryanpujo/blog-app/internal/adapter"
)

type registry struct {
	DB database.DatabaseOperations
}

func New(db database.DatabaseOperations) registry {
	return registry{
		DB: db,
	}
}

func (r registry) NewAppController() adapter.AppController {
	return adapter.AppController{
		UserController:  r.NewUserController(),
		StoryController: r.NewStoryController(),
	}
}

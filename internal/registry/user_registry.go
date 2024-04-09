package registry

import (
	"github.com/ryanpujo/blog-app/internal/user/controllers"
	"github.com/ryanpujo/blog-app/internal/user/repositories"
	"github.com/ryanpujo/blog-app/internal/user/services"
)

func (r registry) NewUserRepository() repositories.UserRepository {
	return repositories.NewUserRepository(r.DB)
}

func (r registry) NewUserService() services.UserService {
	return services.NewUserService(r.NewUserRepository())
}

func (r registry) NewUserController() controllers.UserController {
	return controllers.NewUserController(r.NewUserService())
}

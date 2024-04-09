package registry

import (
	"database/sql"

	"github.com/ryanpujo/blog-app/internal/adapter"
)

type registry struct {
	DB *sql.DB
}

func New(db *sql.DB) registry {
	return registry{
		DB: db,
	}
}

func (r registry) NewAppController() adapter.AppController {
	return adapter.AppController{
		UserController: r.NewUserController(),
	}
}

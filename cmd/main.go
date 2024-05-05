package main

import (
	"github.com/ryanpujo/blog-app/database"
	"github.com/ryanpujo/blog-app/internal/registry"
	"github.com/ryanpujo/blog-app/internal/route"
)

func main() {
	registry := registry.New(database.EstablishDBConnectionWithRetry())
	app := Application(WithPort(4000))
	app.Serve(route.Route(registry.NewAppController()))
}

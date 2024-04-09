package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/blog-app/internal/adapter"
)

var mux = gin.Default()

func Route(app adapter.AppController) *gin.Engine {
	UserRoute(app.UserController)
	return mux
}

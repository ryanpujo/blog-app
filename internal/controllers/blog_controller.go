package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/blog-app/internal/services"
)

type BlogController interface {
	Create(c *gin.Context)
}

type blogController struct {
	s services.BlogService
}

func NewBlogController(s services.BlogService) *blogController {
	return &blogController{
		s: s,
	}
}

func (b *blogController) Create(c *gin.Context) {

}

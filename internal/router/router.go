package router

import (
	"snipqurl/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetUp(h *handler.URLHandler) *gin.Engine {
	r := gin.Default()

	// Serve static files
	r.StaticFile("/", "./static/index.html")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.Static("/static", "./static")

	// API routes
	api := r.Group("/api")
	{
		api.POST("/shorten", h.Shorten)
		api.POST("/qr", h.GenerateQR)
	}

	r.GET("/:code", h.Redirect)

	return r
}

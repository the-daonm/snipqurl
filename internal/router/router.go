package router

import (
	"snipqurl/internal/handler"
	"snipqurl/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetUp(h *handler.URLHandler) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Metrics())

	// Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

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

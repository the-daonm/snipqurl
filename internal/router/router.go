package router

import (
	"snipqurl/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetUp(h *handler.URLHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/shorten", h.Shorten)
		api.POST("/qr", h.GenerateQR)
	}

	r.GET("/:code", h.Redirect)

	return r
}

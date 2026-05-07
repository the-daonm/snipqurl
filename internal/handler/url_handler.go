package handler

import (
	"fmt"
	"net/http"
	"os"

	"snipqurl/internal/service"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	svc service.URLService
}

func New(svc service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

func (h *URLHandler) Shorten(c *gin.Context) {
	type shortenRequest struct {
		URL string `json:"url"`
	}

	var req shortenRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.svc.Shorten(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	domain := os.Getenv("DOMAIN")
	c.JSON(http.StatusOK, gin.H{"code": fmt.Sprintf("https://%s/%s", domain, u.ShortCode)})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	code := c.Param("code")
	u, err := h.svc.GetOriginalURL(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
		return
	}
	c.Redirect(http.StatusMovedPermanently, u.OriginalURL)
}

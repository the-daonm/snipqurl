package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"snipqurl/internal/service"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	svc service.URLService
}

func New(svc service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

type request struct {
	URL       string `json:"url"`
	Alias     string `json:"alias"`
	ExpiresIn string `json:"expires_in"` // e.g. "1h", "24h"
}

func (h *URLHandler) Shorten(c *gin.Context) {
	var req request
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expiresAt *time.Time
	durationStr := req.ExpiresIn
	if durationStr == "" {
		durationStr = "720h" // Default to 30 days
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_in format"})
		return
	}
	t := time.Now().Add(duration)
	expiresAt = &t

	u, err := h.svc.Shorten(req.URL, req.Alias, expiresAt)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) || errors.Is(err, service.ErrAliasTaken) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
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
		if errors.Is(err, service.ErrExpired) {
			c.JSON(http.StatusGone, gin.H{"error": "url has expired"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
		return
	}
	c.Redirect(http.StatusMovedPermanently, u.OriginalURL)
}

func (h *URLHandler) GenerateQR(c *gin.Context) {
	var req request
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	png, err := h.svc.GenerateQR(req.URL)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}

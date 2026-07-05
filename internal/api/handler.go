package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"rate-limiter/internal/manager"
)

type Handler struct {
	manager manager.RateLimiterManager
}

func NewHandler(m manager.RateLimiterManager) *Handler {
	return &Handler{
		manager: m,
	}
}

func (h *Handler) CheckRateLimit(c *gin.Context) {

	var req CheckRequest

	// Parse JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	// Call business layer
	decision, err := h.manager.Check(c.Request.Context(), req.ClientKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
		})
		return
	}

	// Standard Rate Limit Headers
	setRateLimitHeaders(c, decision)

	status := http.StatusOK
	response := CheckResponse{
		Allowed: true,
	}

	if !decision.Allowed {
		status = http.StatusTooManyRequests
		response.Allowed = false
		response.Message = "rate limit exceeded"
	}

	c.JSON(status, response)
}
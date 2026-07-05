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
	c.Header("X-RateLimit-Limit", decision.LimitString())
	c.Header("X-RateLimit-Remaining", decision.RemainingString())
	c.Header("X-RateLimit-Reset", decision.ResetString())

	// Response
	if !decision.Allowed {
		c.Header("Retry-After", decision.RetryAfterString())

		c.JSON(http.StatusTooManyRequests, CheckResponse{
			Allowed: false,
			Message: "rate limit exceeded",
		})
		return
	}

	c.JSON(http.StatusOK, CheckResponse{
		Allowed: true,
	})
}
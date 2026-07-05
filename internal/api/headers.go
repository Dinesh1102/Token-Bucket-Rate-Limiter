package api

import (
	"github.com/gin-gonic/gin"

	"rate-limiter/internal/models"
)

func setRateLimitHeaders(c *gin.Context, d *models.Decision) {
	c.Header("X-RateLimit-Limit", d.LimitString())
	c.Header("X-RateLimit-Remaining", d.RemainingString())
	c.Header("X-RateLimit-Reset", d.ResetString())

	if !d.Allowed {
		c.Header("Retry-After", d.RetryAfterString())
	}
}
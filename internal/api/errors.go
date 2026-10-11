package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"rate-limiter/internal/algorithm"
	"rate-limiter/internal/config"
	"rate-limiter/internal/manager"
)

func writeCheckError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, config.ErrClientNotFound),
		errors.Is(err, manager.ErrClientConfigNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "client not found",
		})

	case errors.Is(err, algorithm.ErrUnknownAlgorithm),
		errors.Is(err, manager.ErrAlgorithmNotSet),
		errors.Is(err, algorithm.ErrInvalidRateLimitConfig):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})

	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
		})
	}
}

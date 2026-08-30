package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"rate-limiter/internal/config"
	"rate-limiter/internal/manager"
	"rate-limiter/internal/models"
)

type Handler struct {
	manager       manager.RateLimiterManager
	configManager config.ConfigManager
}

func NewHandler(m manager.RateLimiterManager, cm config.ConfigManager) *Handler {
	return &Handler{
		manager:       m,
		configManager: cm,
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

func (h *Handler) CreateClient(c *gin.Context) {
	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	clientConfig := models.ClientConfig{
		ClientKey: req.ClientKey,
		Algorithm: models.AlgorithmName(req.Algorithm),
		Rate:      req.Rate,
		Burst:     req.Burst,
	}

	if err := h.configManager.CreateClient(clientConfig); err != nil {
		if errors.Is(err, config.ErrClientAlreadyExists) {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "client already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, clientResponse(clientConfig))
}

func (h *Handler) GetClient(c *gin.Context) {
	clientKey := c.Param("clientKey")

	clientConfig, err := h.configManager.GetClient(clientKey)
	if err != nil {
		if errors.Is(err, config.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "client not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, clientResponse(*clientConfig))
}

func (h *Handler) UpdateClient(c *gin.Context) {
	clientKey := c.Param("clientKey")

	var req UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	clientConfig := models.ClientConfig{
		ClientKey: clientKey,
		Algorithm: models.AlgorithmName(req.Algorithm),
		Rate:      req.Rate,
		Burst:     req.Burst,
	}

	if err := h.configManager.UpdateClient(clientConfig); err != nil {
		if errors.Is(err, config.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "client not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, clientResponse(clientConfig))
}

func (h *Handler) DeleteClient(c *gin.Context) {
	clientKey := c.Param("clientKey")

	if err := h.configManager.DeleteClient(clientKey); err != nil {
		if errors.Is(err, config.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "client not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func clientResponse(config models.ClientConfig) ClientResponse {
	return ClientResponse{
		ClientKey: config.ClientKey,
		Algorithm: config.Algorithm.String(),
		Rate:      config.Rate,
		Burst:     config.Burst,
	}
}

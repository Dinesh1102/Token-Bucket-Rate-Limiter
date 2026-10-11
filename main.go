package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"rate-limiter/internal/algorithm"
	"rate-limiter/internal/api"
	"rate-limiter/internal/config"
	"rate-limiter/internal/manager"
	"rate-limiter/internal/metrics"
	"rate-limiter/internal/state"
)

func main() {
	configManager := config.NewMemoryConfigManager()
	stateManager := state.NewMemoryStateManager()
	algorithmManager := algorithm.NewManager(
		algorithm.NewTokenBucket(),
		algorithm.NewSlidingWindow(),
	)
	metricsCollector := metrics.NewNoopCollector()

	rateLimiterManager := manager.NewRateLimiterManager(
		configManager,
		stateManager,
		algorithmManager,
		metricsCollector,
	)
	handler := api.NewHandler(rateLimiterManager, configManager)

	router := gin.Default()
	api.RegisterRoutes(router, handler)

	addr := serverAddress()
	log.Printf("rate limiter listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func serverAddress() string {
	if addr := os.Getenv("ADDR"); addr != "" {
		return addr
	}
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}

	return ":8080"
}

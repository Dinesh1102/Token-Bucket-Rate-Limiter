package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"rate-limiter/internal/algorithm"
	"rate-limiter/internal/config"
	"rate-limiter/internal/manager"
	"rate-limiter/internal/metrics"
	"rate-limiter/internal/models"
	"rate-limiter/internal/state"
)

func TestAdminClientCRUD(t *testing.T) {
	router, _ := newTestRouter(time.Date(2026, 10, 11, 10, 0, 0, 0, time.UTC))

	createBody := `{"clientKey":"client-1","algorithm":"TOKEN_BUCKET","rate":2,"burst":5}`
	createResponse := performRequest(router, http.MethodPost, "/v1/admin/client", createBody)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d; body = %s", createResponse.Code, http.StatusCreated, createResponse.Body.String())
	}

	getResponse := performRequest(router, http.MethodGet, "/v1/admin/client/client-1", "")
	if getResponse.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d; body = %s", getResponse.Code, http.StatusOK, getResponse.Body.String())
	}

	var client ClientResponse
	if err := json.Unmarshal(getResponse.Body.Bytes(), &client); err != nil {
		t.Fatalf("unmarshal get response: %v", err)
	}
	if client.ClientKey != "client-1" || client.Algorithm != "TOKEN_BUCKET" || client.Rate != 2 || client.Burst != 5 {
		t.Fatalf("client response = %+v, want client-1 TOKEN_BUCKET rate 2 burst 5", client)
	}

	updateBody := `{"algorithm":"TOKEN_BUCKET","rate":10,"burst":20}`
	updateResponse := performRequest(router, http.MethodPut, "/v1/admin/client/client-1", updateBody)
	if updateResponse.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d; body = %s", updateResponse.Code, http.StatusOK, updateResponse.Body.String())
	}

	if err := json.Unmarshal(updateResponse.Body.Bytes(), &client); err != nil {
		t.Fatalf("unmarshal update response: %v", err)
	}
	if client.Rate != 10 || client.Burst != 20 {
		t.Fatalf("updated client = %+v, want rate 10 burst 20", client)
	}

	deleteResponse := performRequest(router, http.MethodDelete, "/v1/admin/client/client-1", "")
	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d; body = %s", deleteResponse.Code, http.StatusNoContent, deleteResponse.Body.String())
	}

	missingResponse := performRequest(router, http.MethodGet, "/v1/admin/client/client-1", "")
	if missingResponse.Code != http.StatusNotFound {
		t.Fatalf("missing get status = %d, want %d; body = %s", missingResponse.Code, http.StatusNotFound, missingResponse.Body.String())
	}
}

func TestAdminClientValidationRejectsInvalidConfig(t *testing.T) {
	router, _ := newTestRouter(time.Date(2026, 10, 11, 10, 0, 0, 0, time.UTC))

	tests := []struct {
		name string
		body string
	}{
		{
			name: "unknown algorithm",
			body: `{"clientKey":"client-1","algorithm":"UNKNOWN","rate":2,"burst":5}`,
		},
		{
			name: "zero rate",
			body: `{"clientKey":"client-1","algorithm":"TOKEN_BUCKET","rate":0,"burst":5}`,
		},
		{
			name: "zero burst",
			body: `{"clientKey":"client-1","algorithm":"TOKEN_BUCKET","rate":2,"burst":0}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := performRequest(router, http.MethodPost, "/v1/admin/client", tt.body)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusBadRequest, response.Body.String())
			}
		})
	}
}

func TestCheckRateLimitAllowsThenDeniesTokenBucket(t *testing.T) {
	router, _ := newTestRouter(time.Date(2026, 10, 11, 10, 0, 0, 0, time.UTC))

	createBody := `{"clientKey":"client-1","algorithm":"TOKEN_BUCKET","rate":1,"burst":2}`
	createResponse := performRequest(router, http.MethodPost, "/v1/admin/client", createBody)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d; body = %s", createResponse.Code, http.StatusCreated, createResponse.Body.String())
	}

	checkBody := `{"clientKey":"client-1"}`
	first := performRequest(router, http.MethodPost, "/v1/check", checkBody)
	if first.Code != http.StatusOK {
		t.Fatalf("first check status = %d, want %d; body = %s", first.Code, http.StatusOK, first.Body.String())
	}
	if first.Header().Get("X-RateLimit-Limit") != "2" || first.Header().Get("X-RateLimit-Remaining") != "1" {
		t.Fatalf("first rate headers limit=%q remaining=%q", first.Header().Get("X-RateLimit-Limit"), first.Header().Get("X-RateLimit-Remaining"))
	}

	second := performRequest(router, http.MethodPost, "/v1/check", checkBody)
	if second.Code != http.StatusOK {
		t.Fatalf("second check status = %d, want %d; body = %s", second.Code, http.StatusOK, second.Body.String())
	}
	if second.Header().Get("X-RateLimit-Remaining") != "0" {
		t.Fatalf("second remaining header = %q, want 0", second.Header().Get("X-RateLimit-Remaining"))
	}

	third := performRequest(router, http.MethodPost, "/v1/check", checkBody)
	if third.Code != http.StatusTooManyRequests {
		t.Fatalf("third check status = %d, want %d; body = %s", third.Code, http.StatusTooManyRequests, third.Body.String())
	}
	if third.Header().Get("Retry-After") != "1" {
		t.Fatalf("retry-after = %q, want 1", third.Header().Get("Retry-After"))
	}
}

func TestCheckRateLimitMapsKnownErrors(t *testing.T) {
	router, cfg := newTestRouter(time.Date(2026, 10, 11, 10, 0, 0, 0, time.UTC))

	missingResponse := performRequest(router, http.MethodPost, "/v1/check", `{"clientKey":"missing"}`)
	if missingResponse.Code != http.StatusNotFound {
		t.Fatalf("missing client status = %d, want %d; body = %s", missingResponse.Code, http.StatusNotFound, missingResponse.Body.String())
	}

	if err := cfg.CreateClient(configForTest("client-1", "SLIDING_WINDOW", 1, 1)); err != nil {
		t.Fatalf("CreateClient() error = %v", err)
	}
	unwiredAlgorithmResponse := performRequest(router, http.MethodPost, "/v1/check", `{"clientKey":"client-1"}`)
	if unwiredAlgorithmResponse.Code != http.StatusBadRequest {
		t.Fatalf("unwired algorithm status = %d, want %d; body = %s", unwiredAlgorithmResponse.Code, http.StatusBadRequest, unwiredAlgorithmResponse.Body.String())
	}

	if err := cfg.CreateClient(configForTest("client-2", "TOKEN_BUCKET", 0, 1)); err != nil {
		t.Fatalf("CreateClient() error = %v", err)
	}
	invalidConfigResponse := performRequest(router, http.MethodPost, "/v1/check", `{"clientKey":"client-2"}`)
	if invalidConfigResponse.Code != http.StatusBadRequest {
		t.Fatalf("invalid config status = %d, want %d; body = %s", invalidConfigResponse.Code, http.StatusBadRequest, invalidConfigResponse.Body.String())
	}
}

func newTestRouter(now time.Time) (*gin.Engine, *config.MemoryConfigManager) {
	gin.SetMode(gin.TestMode)

	configManager := config.NewMemoryConfigManager()
	stateManager := state.NewMemoryStateManager()
	algorithmManager := algorithm.NewManager(
		algorithm.NewTokenBucketWithClock(func() time.Time { return now }),
		nil,
	)
	rateLimiterManager := manager.NewRateLimiterManager(
		configManager,
		stateManager,
		algorithmManager,
		metrics.NewNoopCollector(),
	)
	handler := NewHandler(rateLimiterManager, configManager)

	router := gin.New()
	RegisterRoutes(router, handler)

	return router, configManager
}

func performRequest(router http.Handler, method string, path string, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	return response
}

func configForTest(clientKey string, algorithmName string, rate float64, burst int) models.ClientConfig {
	return models.ClientConfig{
		ClientKey: clientKey,
		Algorithm: models.AlgorithmName(algorithmName),
		Rate:      rate,
		Burst:     burst,
	}
}

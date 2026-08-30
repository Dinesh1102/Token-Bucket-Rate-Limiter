package config

import (
	"errors"
	"testing"

	"rate-limiter/internal/models"
)

func TestMemoryConfigManagerCRUD(t *testing.T) {
	manager := NewMemoryConfigManager()
	clientConfig := models.ClientConfig{
		ClientKey: "client-1",
		Algorithm: models.AlgorithmTokenBucket,
		Rate:      2,
		Burst:     5,
	}

	if err := manager.CreateClient(clientConfig); err != nil {
		t.Fatalf("CreateClient() error = %v", err)
	}

	got, err := manager.GetClient("client-1")
	if err != nil {
		t.Fatalf("GetClient() error = %v", err)
	}
	if *got != clientConfig {
		t.Fatalf("GetClient() = %+v, want %+v", *got, clientConfig)
	}

	updated := clientConfig
	updated.Rate = 10
	if err := manager.UpdateClient(updated); err != nil {
		t.Fatalf("UpdateClient() error = %v", err)
	}

	got, err = manager.GetClient("client-1")
	if err != nil {
		t.Fatalf("GetClient() after update error = %v", err)
	}
	if got.Rate != 10 {
		t.Fatalf("Rate = %v, want 10", got.Rate)
	}

	if err := manager.DeleteClient("client-1"); err != nil {
		t.Fatalf("DeleteClient() error = %v", err)
	}

	_, err = manager.GetClient("client-1")
	if !errors.Is(err, ErrClientNotFound) {
		t.Fatalf("GetClient() error = %v, want ErrClientNotFound", err)
	}
}

func TestMemoryConfigManagerRejectsDuplicateCreate(t *testing.T) {
	manager := NewMemoryConfigManager()
	clientConfig := models.ClientConfig{
		ClientKey: "client-1",
		Algorithm: models.AlgorithmTokenBucket,
		Rate:      2,
		Burst:     5,
	}

	if err := manager.CreateClient(clientConfig); err != nil {
		t.Fatalf("CreateClient() error = %v", err)
	}
	if err := manager.CreateClient(clientConfig); !errors.Is(err, ErrClientAlreadyExists) {
		t.Fatalf("CreateClient() duplicate error = %v, want ErrClientAlreadyExists", err)
	}
}

package state

import (
	"testing"
	"time"

	"rate-limiter/internal/models"
)

func TestMemoryStateManagerLoadMissingReturnsZeroState(t *testing.T) {
	manager := NewMemoryStateManager()

	got, err := manager.Load("client-1")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got != (models.State{}) {
		t.Fatalf("Load() = %+v, want zero state", got)
	}
}

func TestMemoryStateManagerSaveAndLoad(t *testing.T) {
	manager := NewMemoryStateManager()
	want := models.State{
		TokenBucket: models.TokenBucketState{
			Tokens:       3,
			LastRefilled: time.Date(2026, 8, 30, 10, 0, 0, 0, time.UTC),
		},
	}

	if err := manager.Save("client-1", want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := manager.Load("client-1")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got != want {
		t.Fatalf("Load() = %+v, want %+v", got, want)
	}
}

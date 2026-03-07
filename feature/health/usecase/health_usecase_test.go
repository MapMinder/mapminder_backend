package usecase

import (
	"context"
	"testing"
)

type mockHealthRepository struct{}

func (m *mockHealthRepository) CheckSystemHealth() error {
	return nil
}

func TestHealthUsecase_CheckHealth(t *testing.T) {
	mockRepo := &mockHealthRepository{}
	usecase := NewHealthUsecase(mockRepo)

	result, err := usecase.CheckHealth(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Status != "ok" {
		t.Errorf("Expected status 'ok', got %s", result.Status)
	}
}


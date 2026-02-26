package service

import (
	"testing"

	"github.com/weibaohui/k8m/pkg/plugins/modules/harbor/models"
)

func TestNewHarborClient(t *testing.T) {
	registry := &models.HarborRegistry{
		Name:     "test-harbor",
		URL:      "https://harbor.example.com",
		Username: "admin",
		Password: "Harbor12345",
		Insecure: true,
	}

	client := NewHarborClient(registry)

	if client == nil {
		t.Fatal("NewHarborClient returned nil")
	}

	if client.BaseURL != "https://harbor.example.com" {
		t.Errorf("Expected BaseURL to be 'https://harbor.example.com', got '%s'", client.BaseURL)
	}

	if client.Username != "admin" {
		t.Errorf("Expected Username to be 'admin', got '%s'", client.Username)
	}

	if client.Password != "Harbor12345" {
		t.Errorf("Expected Password to be 'Harbor12345', got '%s'", client.Password)
	}

	if client.HTTPClient == nil {
		t.Fatal("HTTPClient is nil")
	}
}

func TestHarborClientURLTrimming(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "URL with trailing slash",
			inputURL: "https://harbor.example.com/",
			expected: "https://harbor.example.com",
		},
		{
			name:     "URL without trailing slash",
			inputURL: "https://harbor.example.com",
			expected: "https://harbor.example.com",
		},
		{
			name:     "URL with multiple trailing slashes",
			inputURL: "https://harbor.example.com///",
			expected: "https://harbor.example.com//",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := &models.HarborRegistry{
				URL:      tt.inputURL,
				Username: "admin",
				Password: "password",
			}

			client := NewHarborClient(registry)

			if client.BaseURL != tt.expected {
				t.Errorf("Expected BaseURL to be '%s', got '%s'", tt.expected, client.BaseURL)
			}
		})
	}
}

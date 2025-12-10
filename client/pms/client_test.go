/*
 * Copyright (C) 2025 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package pms

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"intel/tac/v1/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewPmsClient(t *testing.T) {
	// Arrange
	baseURL, _ := url.Parse("https://example.com")
	apiKey := "test-api-key"
	client := &http.Client{}

	// Act
	pmsClient := NewPmsClient(client, baseURL, apiKey)
	// Assert
	assert.NotNil(t, pmsClient)
}

func TestCreatePolicy(t *testing.T) {
	// Arrange
	policyID := uuid.New()
	policyName := "test-policy"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.PolicyResponse{
			CommonPolicy: models.CommonPolicy{
				PolicyId:   policyID,
				PolicyName: policyName,
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	request := &models.PolicyRequest{
		CommonPolicy: models.CommonPolicy{
			PolicyName: policyName,
		},
	}
	// Act
	result, err := client.CreatePolicy(request)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestDeletePolicy(t *testing.T) {
	// Arrange
	policyID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	err := client.DeletePolicy(policyID)
	// Assert
	assert.NoError(t, err)
}

func TestGetPolicy(t *testing.T) {
	// Arrange
	policyID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.PolicyResponse{
			CommonPolicy: models.CommonPolicy{
				PolicyId: policyID,
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	result, err := client.GetPolicy(policyID)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSearchPolicy(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []models.PolicyResponse{
			{CommonPolicy: models.CommonPolicy{PolicyId: uuid.New()}},
			{CommonPolicy: models.CommonPolicy{PolicyId: uuid.New()}},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	result, err := client.SearchPolicy()
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
}

func TestUpdatePolicy(t *testing.T) {
	// Arrange
	policyID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.PolicyResponse{
			CommonPolicy: models.CommonPolicy{
				PolicyId: policyID,
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.PolicyUpdateRequest{
		PolicyId: policyID,
	}
	// Act
	result, err := client.UpdatePolicy(request)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCreatePolicyError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid policy"}`))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.PolicyRequest{
		CommonPolicy: models.CommonPolicy{
			PolicyName: "test",
		},
	}
	// Act
	_, err := client.CreatePolicy(request)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid policy")
}

func TestDeletePolicyError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	err := client.DeletePolicy(uuid.New())
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "404 Not Found")
}

func TestGetPolicyError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	_, err := client.GetPolicy(uuid.New())
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "404 Not Found")
}

func TestSearchPolicyError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	_, err := client.SearchPolicy()
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "500 Internal Server Error")
}

func TestUpdatePolicyError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.PolicyUpdateRequest{
		PolicyId: uuid.New(),
	}
	// Act
	_, err := client.UpdatePolicy(request)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "400 Bad Request")
}

// Invalid JSON tests to improve coverage

func TestCreatePolicyInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{invalid json"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.PolicyRequest{
		CommonPolicy: models.CommonPolicy{
			PolicyName: "test",
		},
	}
	// Act
	_, err := client.CreatePolicy(request)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error unmarshalling response")
}

func TestGetPolicyInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not valid json"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	_, err := client.GetPolicy(uuid.New())
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error unmarshalling policy response")
}

func TestSearchPolicyInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[incomplete"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	_, err := client.SearchPolicy()
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error unmarshalling response")
}

func TestUpdatePolicyInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("}{"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.PolicyUpdateRequest{
		PolicyId: uuid.New(),
	}
	// Act
	_, err := client.UpdatePolicy(request)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error unmarshalling response")
}

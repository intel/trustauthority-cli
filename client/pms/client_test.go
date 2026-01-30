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

// RIM Tests

func TestCreateRim(t *testing.T) {
	// Arrange
	rimID := uuid.New()
	rimName := "test.rim.name"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.SignedRimResponse{
			CommonRim: models.CommonRim{
				Id:   rimID,
				Name: rimName,
			},
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	request := &models.RimCreateRequest{
		Name:    rimName,
		Content: json.RawMessage(`{"measurements":[]}`),
	}
	// Act
	result, err := client.CreateRim(request)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, rimID, result.Id)
	assert.Equal(t, rimName, result.Name)
}

func TestCreateRimError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid rim"}`))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.RimCreateRequest{
		Name: "test",
	}
	// Act
	_, err := client.CreateRim(request)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid rim")
}

func TestCreateRimInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{invalid json"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.RimCreateRequest{
		Name: "test",
	}
	// Act
	_, err := client.CreateRim(request)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error unmarshalling response")
}

func TestDeleteRim(t *testing.T) {
	// Arrange
	rimID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	err := client.DeleteRim(rimID)
	// Assert
	assert.NoError(t, err)
}

func TestDeleteRimError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	err := client.DeleteRim(uuid.New())
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "404 Not Found")
}

func TestGetRim(t *testing.T) {
	// Arrange
	rimID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.SignedRimResponse{
			CommonRim: models.CommonRim{Id: rimID},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	result, err := client.GetRim(rimID)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, rimID, result.Id)
}

func TestGetRimError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	_, err := client.GetRim(uuid.New())
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "404 Not Found")
}

func TestGetRimInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not valid json"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	_, err := client.GetRim(uuid.New())
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error unmarshalling rim response")
}

func TestSearchRim(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify rimName query parameter is passed
		name := r.URL.Query().Get("rimName")
		assert.Equal(t, "test.rim", name)

		response := []models.SignedRimResponse{
			{CommonRim: models.CommonRim{Id: uuid.New()}},
			{CommonRim: models.CommonRim{Id: uuid.New()}},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	queryParams := map[string]string{"rimName": "test.rim"}
	result, err := client.SearchRim(queryParams)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
}

func TestSearchRimNoName(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify no query parameters when empty map
		assert.Equal(t, 0, len(r.URL.Query()))

		response := []models.SignedRimResponse{
			{CommonRim: models.CommonRim{Id: uuid.New()}},
			{CommonRim: models.CommonRim{Id: uuid.New()}},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	queryParams := make(map[string]string)
	result, err := client.SearchRim(queryParams)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
}

func TestSearchRimWithMultipleParams(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify multiple query parameters are passed
		query := r.URL.Query()
		assert.Equal(t, "test.rim", query.Get("rimName"))
		assert.Equal(t, "true", query.Get("includeOwnPrivate"))
		assert.Equal(t, "true", query.Get("excludeContent"))

		response := []models.SignedRimResponse{
			{CommonRim: models.CommonRim{Id: uuid.New()}},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	queryParams := map[string]string{
		"rimName":           "test.rim",
		"includeOwnPrivate": "true",
		"excludeContent":    "true",
	}
	result, err := client.SearchRim(queryParams)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
}

func TestSearchRimError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	queryParams := make(map[string]string)
	_, err := client.SearchRim(queryParams)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "500 Internal Server Error")
}

func TestSearchRimInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[incomplete"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	queryParams := make(map[string]string)
	_, err := client.SearchRim(queryParams)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error unmarshalling response")
}

func TestUpdateRim(t *testing.T) {
	// Arrange
	rimID := uuid.New()
	rimName := "updated.rim"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		response := models.SignedRimResponse{
			CommonRim: models.CommonRim{
				Id:   rimID,
				Name: rimName,
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	request := &models.RimUpdateRequest{
		Id:          rimID,
		Content:     json.RawMessage(`{"measurements":[]}`),
		Description: "Updated description",
	}
	// Act
	result, err := client.UpdateRim(request)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, rimID, result.Id)
}

func TestUpdateRimError(t *testing.T) {
	// Arrange
	rimID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid update"}`))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	request := &models.RimUpdateRequest{
		Id: rimID,
	}
	// Act
	_, err := client.UpdateRim(request)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid update")
}

func TestUpdateRimInvalidJSON(t *testing.T) {
	// Arrange
	rimID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{invalid json"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewPmsClient(http.DefaultClient, baseURL, "test-key")
	request := &models.RimUpdateRequest{
		Id: rimID,
	}
	// Act
	_, err := client.UpdateRim(request)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error unmarshalling response")
}

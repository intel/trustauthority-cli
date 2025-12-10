/*
 * Copyright (C) 2025 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package tms

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

func TestNewTmsClient(t *testing.T) {
	// Arrange
	baseURL, _ := url.Parse("https://example.com")
	apiKey := "test-api-key"
	client := &http.Client{}
	// Act
	tmsClient := NewTmsClient(client, baseURL, apiKey)
	// Assert
	assert.NotNil(t, tmsClient)
}

func TestCreateApiClient(t *testing.T) {
	// Arrange
	serviceID := uuid.New()
	apiClientID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.ApiClientDetail{
			ID: apiClientID,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.CreateApiClient{
		ServiceId: serviceID,
	}
	// Act
	result, err := client.CreateApiClient(request)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetApiClient(t *testing.T) {
	// Arrange
	serviceID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []models.ApiClient{
			{ID: uuid.New()},
			{ID: uuid.New()},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	// Act
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	result, err := client.GetApiClient(serviceID)
	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDeleteApiClient(t *testing.T) {
	// Arrange
	serviceID := uuid.New()
	apiClientID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	// Act
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	err := client.DeleteApiClient(serviceID, apiClientID)
	// Assert
	assert.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	// Arrange
	userID := uuid.New()
	email := "test@example.com"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.TenantUser{
			ID:    userID,
			Email: email,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	// Act
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.CreateTenantUser{
		Email: email,
	}
	result, err := client.CreateUser(request)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
func TestGetUsers(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []models.TenantUser{
			{ID: uuid.New(), Email: "user1@example.com"},
			{ID: uuid.New(), Email: "user2@example.com"},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	result, err := client.GetUsers()
	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDeleteUser(t *testing.T) {
	// Arrange
	userID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	// Act
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	err := client.DeleteUser(userID)
	// Assert
	assert.NoError(t, err)
}

func TestGetServices(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []models.Service{
			{ID: uuid.New()},
			{ID: uuid.New()},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	// Act
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	result, err := client.GetServices()
	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}
func TestRetrieveService(t *testing.T) {
	// Arrange
	serviceID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.ServiceDetail{
			ID: serviceID,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	// Act
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	result, err := client.RetrieveService(serviceID)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetProducts(t *testing.T) {
	// Arrange
	serviceOfferID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []models.Product{
			{ID: uuid.New()},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	// Act
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	result, err := client.GetProducts(serviceOfferID)
	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGetServiceOffers(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []models.ServiceOffer{
			{ID: uuid.New(), Name: "Offer 1"},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	result, err := client.GetServiceOffers()
	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestCreateTenantTag(t *testing.T) {
	// Arrange
	tagID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.Tag{
			ID:   &tagID,
			Name: "test-tag",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.TagCreate{
		Name: "test-tag",
	}
	// Act
	result, err := client.CreateTenantTag(request)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetTenantTags(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.Tags{
			Tags: []models.Tag{},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	result, err := client.GetTenantTags()
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
func TestDeleteTenantTag(t *testing.T) {
	// Arrange
	tagID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	err := client.DeleteTenantTag(tagID)
	// Assert
	assert.NoError(t, err)
}

func TestGetPlans(t *testing.T) {
	// Arrange
	serviceOfferID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []models.Plan{
			{ID: uuid.New()},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	result, err := client.GetPlans(serviceOfferID)
	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}
func TestRetrievePlan(t *testing.T) {
	// Arrange
	serviceOfferID := uuid.New()
	planID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.PlanProducts{
			ID: planID,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	result, err := client.RetrievePlan(serviceOfferID, planID)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUpdateTenantSettings(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.AttestationFailureEmail{
			AttestationFailureEmail: "admin@example.com",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.AttestationFailureEmail{
		AttestationFailureEmail: "admin@example.com",
	}
	// Act
	result, err := client.UpdateTenantSettings(request)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
func TestGetTenantSettings(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.AttestationFailureEmail{
			AttestationFailureEmail: "admin@example.com",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	result, err := client.GetTenantSettings()
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUpdateApiClient(t *testing.T) {
	// Arrange
	serviceID := uuid.New()
	apiClientID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.ApiClient{
			ID: apiClientID,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.UpdateApiClient{
		ServiceId: serviceID,
	}
	// Act
	result, err := client.UpdateApiClient(request, apiClientID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestRetrieveApiClient(t *testing.T) {
	// Arrange
	serviceID := uuid.New()
	apiClientID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.ApiClientDetail{
			ID: apiClientID,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	result, err := client.RetrieveApiClient(serviceID, apiClientID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetApiClientPolicies(t *testing.T) {
	// Arrange
	serviceID := uuid.New()
	apiClientID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.ApiClientPolicies{
			PolicyIds: []uuid.UUID{},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	result, err := client.GetApiClientPolicies(serviceID, apiClientID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetApiClientTagValues(t *testing.T) {
	// Arrange
	serviceID := uuid.New()
	apiClientID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.ApiClientTags{
			TagsValues: []models.ApiClientTagValue{},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	result, err := client.GetApiClientTagValues(serviceID, apiClientID)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUpdateTenantUserRole(t *testing.T) {
	// Arrange
	userID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.TenantUser{
			ID: userID,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.UpdateTenantUserRoles{
		UserId: userID,
	}
	// Act
	result, err := client.UpdateTenantUserRole(request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

// Error handling tests
func TestCreateApiClientError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.CreateApiClient{
		ServiceId: uuid.New(),
	}

	// Act
	_, err := client.CreateApiClient(request)

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "400 Bad Request")
}

func TestDeleteApiClientError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	err := client.DeleteApiClient(uuid.New(), uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "404 Not Found")
}

func TestCreateUserError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.CreateTenantUser{
		Email: "test@example.com",
	}

	// Act
	_, err := client.CreateUser(request)

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "400 Bad Request")
}

func TestDeleteUserError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	err := client.DeleteUser(uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "404 Not Found")
}

func TestDeleteTenantTagError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")
	// Act
	err := client.DeleteTenantTag(uuid.New())
	// Assert
	assert.NotNil(t, err)
}

func TestUpdateApiClientError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.UpdateApiClient{
		ServiceId: uuid.New(),
	}
	// Act
	_, err := client.UpdateApiClient(request, uuid.New())
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "400 Bad Request")
}

func TestCreateTenantTagError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.TagCreate{
		Name: "test-tag",
	}
	// Act
	_, err := client.CreateTenantTag(request)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "400 Bad Request")
}

func TestUpdateTenantSettingsError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.AttestationFailureEmail{
		AttestationFailureEmail: "test@example.com",
	}
	// Act
	_, err := client.UpdateTenantSettings(request)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "400 Bad Request")
}

func TestUpdateTenantUserRoleError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	request := &models.UpdateTenantUserRoles{
		UserId: uuid.New(),
	}
	// Act
	_, err := client.UpdateTenantUserRole(request)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "400 Bad Request")
}

func TestGetApiClientError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetApiClient(uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "500 Internal Server Error")
}

func TestGetApiClientInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetApiClient(uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestRetrieveApiClientError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.RetrieveApiClient(uuid.New(), uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "404 Not Found")
}

func TestRetrieveApiClientInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.RetrieveApiClient(uuid.New(), uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestGetApiClientPoliciesError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetApiClientPolicies(uuid.New(), uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "500 Internal Server Error")
}

func TestGetApiClientPoliciesInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{invalid}"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetApiClientPolicies(uuid.New(), uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestGetApiClientTagValuesError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetApiClientTagValues(uuid.New(), uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "403 Forbidden")
}

func TestGetApiClientTagValuesInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("bad json"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetApiClientTagValues(uuid.New(), uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestGetUsersError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetUsers()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "401 Unauthorized")
}

func TestGetUsersInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[invalid"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetUsers()

	// Assert: Should return JSON parsing error
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestGetServicesError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetServices()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "503 Service Unavailable")
}

func TestGetServicesInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetServices()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestRetrieveServiceError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.RetrieveService(uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "404 Not Found")
}

func TestRetrieveServiceInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{bad"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.RetrieveService(uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestGetProductsError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetProducts(uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "500 Internal Server Error")
}

func TestGetProductsInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[}"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetProducts(uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestGetServiceOffersError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetServiceOffers()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "400 Bad Request")
}

func TestGetServiceOffersInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[invalid"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetServiceOffers()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestGetTenantTagsError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetTenantTags()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "401 Unauthorized")
}

func TestGetTenantTagsInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{incomplete"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetTenantTags()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestGetPlansError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetPlans(uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "403 Forbidden")
}

func TestGetPlansInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[[]]"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetPlans(uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error unmarshalling response")
}

func TestRetrievePlanError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.RetrievePlan(uuid.New(), uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "404 Not Found")
}

func TestRetrievePlanInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("}{"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.RetrievePlan(uuid.New(), uuid.New())

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestGetTenantSettingsError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetTenantSettings()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "500 Internal Server Error")
}

func TestGetTenantSettingsInvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("broken json"))
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client := NewTmsClient(http.DefaultClient, baseURL, "test-key")

	// Act
	_, err := client.GetTenantSettings()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

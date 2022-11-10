/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package tms

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"intel/amber/tac/v1/client"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/models"
	"net/http"
	"net/url"
)

type TmsClient interface {
	CreateApiClient(request *models.CreateApiClient) (*models.ApiClientDetail, error)
	UpdateApiClient(request *models.UpdateApiClient, apiClientid uuid.UUID) (*models.ApiClient, error)
	GetApiClient(serviceId uuid.UUID) ([]models.ApiClient, error)
	RetrieveApiClient(serviceId uuid.UUID, apiClientId uuid.UUID) (*models.ApiClientDetail, error)
	GetApiClientPolicies(serviceId, apiClientId uuid.UUID) (*models.ApiClientPolicies, error)
	GetApiClientTagValues(serviceId, apiClientId uuid.UUID) (*models.ApiClientTagsValues, error)
	DeleteApiClient(serviceId, apiClientId uuid.UUID) error

	CreateService(request *models.CreateService) (*models.Service, error)
	UpdateService(request *models.UpdateService) (*models.Service, error)
	GetServices() ([]models.Service, error)
	RetrieveService(serviceId uuid.UUID) (*models.ServiceDetail, error)
	DeleteService(serviceId uuid.UUID) error

	GetProducts(serviceOfferId uuid.UUID) ([]models.Product, error)

	GetServiceOffers() ([]models.ServiceOffer, error)

	CreateUser(user *models.CreateTenantUser) (*models.TenantUser, error)
	UpdateTenantUserRole(user *models.UpdateTenantUserRoles) (*models.TenantUser, error)
	GetUsers() ([]models.TenantUser, error)
	DeleteUser(userId uuid.UUID) error

	CreateTenantTag(request *models.Tag) (*models.Tag, error)
	GetTenantTags() (*models.Tags, error)
}

//Client Details for TMS client
type tmsClient struct {
	Client   *http.Client
	BaseURL  *url.URL
	TenantId uuid.UUID
	ApiKey   string
}

func NewTmsClient(client *http.Client, qvsURL *url.URL, tenantId uuid.UUID, apiKey string) TmsClient {
	return &tmsClient{
		Client:   client,
		BaseURL:  qvsURL,
		TenantId: tenantId,
		ApiKey:   apiKey,
	}
}

func (pc tmsClient) CreateApiClient(request *models.CreateApiClient) (*models.ApiClientDetail, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling request")
	}

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceApiEndpoint + "/" +
		request.ServiceId.String() + constants.ApiClientResourceEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPost, reqURL.String(), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyCreatedBy, pc.TenantId.String())

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body")
	}

	// Parse response for validation
	var apiClientDetail models.ApiClientDetail
	err = json.Unmarshal(response, &apiClientDetail)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &apiClientDetail, nil
}

func (pc tmsClient) UpdateApiClient(request *models.UpdateApiClient, apiClientId uuid.UUID) (*models.ApiClient, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling request")
	}

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceApiEndpoint + "/" +
		request.ServiceId.String() + constants.ApiClientResourceEndpoint + "/" + apiClientId.String())
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPut, reqURL.String(), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyUpdatedBy, pc.TenantId.String())

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body")
	}

	// Parse response for validation
	var apiClientDetail models.ApiClient
	err = json.Unmarshal(response, &apiClientDetail)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &apiClientDetail, nil
}

func (pc tmsClient) GetApiClient(serviceId uuid.UUID) ([]models.ApiClient, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceApiEndpoint + "/" +
		serviceId.String() + constants.ApiClientResourceEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body")
	}

	// Parse response for validation
	var apiClients []models.ApiClient
	err = json.Unmarshal(response, &apiClients)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return apiClients, nil
}

func (pc tmsClient) RetrieveApiClient(serviceId uuid.UUID, apiClientId uuid.UUID) (*models.ApiClientDetail, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceApiEndpoint + "/" +
		serviceId.String() + constants.ApiClientResourceEndpoint + "/" + apiClientId.String())
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var apiClients models.ApiClientDetail
	err = json.Unmarshal(response, &apiClients)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &apiClients, nil
}

func (pc tmsClient) GetApiClientPolicies(serviceId, apiClientId uuid.UUID) (*models.ApiClientPolicies, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceApiEndpoint + "/" +
		serviceId.String() + constants.ApiClientResourceEndpoint + "/" + apiClientId.String() + constants.PolicyApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var apiClientPolicies models.ApiClientPolicies
	err = json.Unmarshal(response, &apiClientPolicies)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &apiClientPolicies, nil
}

func (pc tmsClient) GetApiClientTagValues(serviceId, apiClientId uuid.UUID) (*models.ApiClientTagsValues, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceApiEndpoint + "/" +
		serviceId.String() + constants.ApiClientResourceEndpoint + "/" + apiClientId.String() + constants.TagApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var apiClientTagsValues models.ApiClientTagsValues
	err = json.Unmarshal(response, &apiClientTagsValues)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &apiClientTagsValues, nil
}

func (pc tmsClient) DeleteApiClient(serviceId, apiClientId uuid.UUID) error {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceApiEndpoint + "/" +
		serviceId.String() + constants.ApiClientResourceEndpoint + "/" + apiClientId.String())
	if err != nil {
		return errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodDelete, reqURL.String(), nil)
	if err != nil {
		return errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyUpdatedBy, pc.TenantId.String())

	_, err = client.SendRequest(pc.Client, req)
	if err != nil {
		return errors.Wrap(err, "Error in response body")
	}

	return nil
}

func (pc tmsClient) CreateUser(user *models.CreateTenantUser) (*models.TenantUser, error) {
	reqBytes, err := json.Marshal(user)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling request")
	}

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.UserApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPost, reqURL.String(), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyCreatedBy, pc.TenantId.String())

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var createUserRes models.TenantUser
	err = json.Unmarshal(response, &createUserRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &createUserRes, nil
}

func (pc tmsClient) UpdateTenantUserRole(request *models.UpdateTenantUserRoles) (*models.TenantUser, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling request")
	}

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.UserApiEndpoint +
		"/" + request.UserId.String())
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPut, reqURL.String(), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyUpdatedBy, pc.TenantId.String())

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var updateUserRes models.TenantUser
	err = json.Unmarshal(response, &updateUserRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &updateUserRes, nil
}

func (pc tmsClient) GetUsers() ([]models.TenantUser, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.UserApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var searchUserRes []models.TenantUser
	err = json.Unmarshal(response, &searchUserRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return searchUserRes, nil
}

func (pc tmsClient) DeleteUser(userId uuid.UUID) error {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.UserApiEndpoint + "/" + userId.String())
	if err != nil {
		return errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodDelete, reqURL.String(), nil)
	if err != nil {
		return errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyUpdatedBy, pc.TenantId.String())

	_, err = client.SendRequest(pc.Client, req)
	if err != nil {
		return errors.Wrap(err, "Error reading response body")
	}

	return nil
}

func (pc tmsClient) CreateService(request *models.CreateService) (*models.Service, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling request")
	}

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPost, reqURL.String(), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyCreatedBy, pc.TenantId.String())

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var serviceDetail models.Service
	err = json.Unmarshal(response, &serviceDetail)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &serviceDetail, nil
}

func (pc tmsClient) UpdateService(request *models.UpdateService) (*models.Service, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling request")
	}

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceApiEndpoint + "/" + request.Id.String())
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPut, reqURL.String(), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyUpdatedBy, pc.TenantId.String())

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var serviceDetail models.Service
	err = json.Unmarshal(response, &serviceDetail)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &serviceDetail, nil
}

func (pc tmsClient) DeleteService(serviceId uuid.UUID) error {

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceApiEndpoint + "/" + serviceId.String())
	if err != nil {
		return errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodDelete, reqURL.String(), nil)
	if err != nil {
		return errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)

	_, err = client.SendRequest(pc.Client, req)
	if err != nil {
		return errors.Wrap(err, "Error in response body")
	}

	if err != nil {
		return errors.Wrap(err, "Error unmarshalling response")
	}
	return nil
}

func (pc tmsClient) GetServices() ([]models.Service, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var searchServiceRes []models.Service
	err = json.Unmarshal(response, &searchServiceRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return searchServiceRes, nil
}

func (pc tmsClient) RetrieveService(id uuid.UUID) (*models.ServiceDetail, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceApiEndpoint + "/" + id.String())
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var retrieveServiceRes *models.ServiceDetail
	err = json.Unmarshal(response, &retrieveServiceRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return retrieveServiceRes, nil
}

func (pc tmsClient) GetProducts(serviceOfferId uuid.UUID) ([]models.Product, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceOfferApiEndpoint + "/" + serviceOfferId.String() + constants.ProductApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var searchServiceRes []models.Product
	err = json.Unmarshal(response, &searchServiceRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return searchServiceRes, nil
}

func (pc tmsClient) GetServiceOffers() ([]models.ServiceOffer, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceOfferApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var searchServiceRes []models.ServiceOffer
	err = json.Unmarshal(response, &searchServiceRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return searchServiceRes, nil
}

func (pc tmsClient) CreateTenantTag(request *models.Tag) (*models.Tag, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TagApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling request")
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPost, reqURL.String(), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyCreatedBy, pc.TenantId.String())

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	log.Info(string(response))

	// Parse response for validation
	var createTagRes models.Tag
	err = json.Unmarshal(response, &createTagRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &createTagRes, nil
}

func (pc tmsClient) GetTenantTags() (*models.Tags, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TagApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var getTagsRes models.Tags
	err = json.Unmarshal(response, &getTagsRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &getTagsRes, nil
}

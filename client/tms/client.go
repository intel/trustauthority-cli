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
	"intel/amber/tac/v1/client"
	"intel/amber/tac/v1/constants"
	models2 "intel/amber/tac/v1/internal/models"
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
	GetApiClientTagValues(serviceId, apiClientId uuid.UUID) (*models.ApiClientTags, error)
	DeleteApiClient(serviceId, apiClientId uuid.UUID) error

	GetServices() ([]models.Service, error)
	RetrieveService(serviceId uuid.UUID) (*models.ServiceDetail, error)

	GetProducts(serviceOfferId uuid.UUID) ([]models.Product, error)

	GetServiceOffers() ([]models.ServiceOffer, error)

	CreateUser(user *models.CreateTenantUser) (*models.TenantUser, error)
	UpdateTenantUserRole(user *models.UpdateTenantUserRoles) (*models.TenantUser, error)
	GetUsers() ([]models.TenantUser, error)
	DeleteUser(userId uuid.UUID) error

	CreateTenantTag(request *models.TagCreate) (*models.Tag, error)
	GetTenantTags() (*models.Tags, error)
	DeleteTenantTag(tagId uuid.UUID) error

	GetPlans(serviceOfferId uuid.UUID) ([]models.Plan, error)
	RetrievePlan(serviceOfferId, planId uuid.UUID) (*models.PlanProducts, error)

	UpdateTenantSettings(settings *models.AttestationFailureEmail) (*models.AttestationFailureEmail, error)
	GetTenantSettings() (*models.AttestationFailureEmail, error)
}

// Client Details for TMS client
type tmsClient struct {
	Client  *http.Client
	BaseURL *url.URL
	ApiKey  string
}

func NewTmsClient(client *http.Client, tmsURL *url.URL, apiKey string) TmsClient {
	return &tmsClient{
		Client:  client,
		BaseURL: tmsURL,
		ApiKey:  apiKey,
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
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body")
	}

	// Parse response for validation
	var apiClientDetail models.ApiClientDetail
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&apiClientDetail)
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
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body")
	}

	// Parse response for validation
	var apiClientDetail models.ApiClient
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&apiClientDetail)
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
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body")
	}

	// Parse response for validation
	var apiClients []models.ApiClient
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&apiClients)
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
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var apiClients models.ApiClientDetail
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&apiClients)
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
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var apiClientPolicies models.ApiClientPolicies
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&apiClientPolicies)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &apiClientPolicies, nil
}

func (pc tmsClient) GetApiClientTagValues(serviceId, apiClientId uuid.UUID) (*models.ApiClientTags, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceApiEndpoint + "/" +
		serviceId.String() + constants.ApiClientResourceEndpoint + "/" + apiClientId.String() + constants.TagApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var apiClientTagsValues models.ApiClientTags
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&apiClientTagsValues)
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
		return errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	_, err = client.SendRequest(pc.Client, req)
	if err != nil {
		return errors.Wrap(err, "Error in response")
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
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var createUserRes models.TenantUser
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&createUserRes)
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
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var updateUserRes models.TenantUser
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&updateUserRes)
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
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var searchUserRes []models.TenantUser
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&searchUserRes)
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
		return errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	_, err = client.SendRequest(pc.Client, req)
	if err != nil {
		return errors.Wrap(err, "Error reading response body")
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
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var searchServiceRes []models.Service
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&searchServiceRes)
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
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var retrieveServiceRes *models.ServiceDetail
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&retrieveServiceRes)
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
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var searchProductsRes []models.Product
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&searchProductsRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return searchProductsRes, nil
}

func (pc tmsClient) GetServiceOffers() ([]models.ServiceOffer, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceOfferApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var searchServiceOfferRes []models.ServiceOffer
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&searchServiceOfferRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return searchServiceOfferRes, nil
}

func (pc tmsClient) CreateTenantTag(request *models.TagCreate) (*models.Tag, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TagApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "Error marshalling request")
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPost, reqURL.String(), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var createTagRes models.Tag
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&createTagRes)
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
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var getTagsRes models.Tags
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&getTagsRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &getTagsRes, nil
}

func (pc tmsClient) DeleteTenantTag(tagId uuid.UUID) error {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TagApiEndpoint + "/" + tagId.String())
	if err != nil {
		return errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodDelete, reqURL.String(), nil)
	if err != nil {
		return errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	_, err = client.SendRequest(pc.Client, req)
	if err != nil {
		return errors.Wrap(err, "Error reading response body")
	}

	return nil
}

func (pc tmsClient) GetPlans(serviceOfferId uuid.UUID) ([]models.Plan, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceOfferApiEndpoint + "/" + serviceOfferId.String() +
		constants.PlanApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var searchPlanRes []models.Plan
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&searchPlanRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return searchPlanRes, nil
}

func (pc tmsClient) RetrievePlan(serviceOfferId, planId uuid.UUID) (*models.PlanProducts, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.ServiceOfferApiEndpoint + "/" + serviceOfferId.String() +
		constants.PlanApiEndpoint + "/" + planId.String())
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var retrievePlanRes models.PlanProducts
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&retrievePlanRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &retrievePlanRes, nil
}

func (pc tmsClient) UpdateTenantSettings(request *models.AttestationFailureEmail) (*models.AttestationFailureEmail, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantsApiEndpoint + constants.SettingsEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling request")
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPut, reqURL.String(), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var notificationUpdateRes models.AttestationFailureEmail
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&notificationUpdateRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &notificationUpdateRes, nil
}

func (pc tmsClient) GetTenantSettings() (*models.AttestationFailureEmail, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantsApiEndpoint + constants.SettingsEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response for validation
	var notificationFetchRes models.AttestationFailureEmail
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&notificationFetchRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &notificationFetchRes, nil
}

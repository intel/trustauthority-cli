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
	CreateSubscription(request *models.CreateSubscription) (*models.SubscriptionDetail, error)
	UpdateSubscription(request *models.UpdateSubscription, subscriptionid uuid.UUID) (*models.Subscription, error)
	GetSubscriptions(serviceId uuid.UUID) ([]models.Subscription, error)
	RetrieveSubscription(serviceId uuid.UUID, subscriptionId uuid.UUID) (*models.SubscriptionDetail, error)
	GetSubscriptionPolicies(serviceId, subscriptionId uuid.UUID) (*models.SubscriptionPolicies, error)
	GetSubscriptionTagValues(serviceId, subscriptionId uuid.UUID) (*models.SubscriptionTagsValues, error)
	DeleteSubscription(serviceId, subscriptionId uuid.UUID) error

	CreateService(request *models.CreateService) (*models.Service, error)
	UpdateService(request *models.UpdateService) (*models.Service, error)
	GetServices() ([]models.Service, error)
	RetrieveService(serviceId uuid.UUID) (*models.ServiceDetail, error)
	DeleteService(serviceId uuid.UUID) error

	GetProducts(serviceOfferId uuid.UUID) ([]models.Product, error)

	GetServiceOffers() ([]models.ServiceOffer, error)

	CreateUser(user *models.CreateTenantUser) (*models.User, error)
	UpdateTenantUserRole(user *models.UpdateTenantUserRoles) (*models.User, error)
	UpdateUser(user *models.UpdateUser) (*models.User, error)
	GetUsers() ([]models.User, error)
	RetrieveUser(userId uuid.UUID) (*models.User, error)
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

func (pc tmsClient) CreateSubscription(request *models.CreateSubscription) (*models.SubscriptionDetail, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling request")
	}

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.ServiceApiEndpoint + "/" +
		request.ServiceId.String() + constants.SubscriptionApiEndpoint)
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
	var subscriptionDetail models.SubscriptionDetail
	err = json.Unmarshal(response, &subscriptionDetail)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &subscriptionDetail, nil
}

func (pc tmsClient) UpdateSubscription(request *models.UpdateSubscription, subscriptionId uuid.UUID) (*models.Subscription, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling request")
	}

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.ServiceApiEndpoint + "/" +
		request.ServiceId.String() + constants.SubscriptionApiEndpoint + "/" + subscriptionId.String())
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
	var subscriptionDetail models.Subscription
	err = json.Unmarshal(response, &subscriptionDetail)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &subscriptionDetail, nil
}

func (pc tmsClient) GetSubscriptions(serviceId uuid.UUID) ([]models.Subscription, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.ServiceApiEndpoint + "/" +
		serviceId.String() + constants.SubscriptionApiEndpoint)
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
	var subscriptions []models.Subscription
	err = json.Unmarshal(response, &subscriptions)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return subscriptions, nil
}

func (pc tmsClient) RetrieveSubscription(serviceId uuid.UUID, subscriptionId uuid.UUID) (*models.SubscriptionDetail, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.ServiceApiEndpoint + "/" +
		serviceId.String() + constants.SubscriptionApiEndpoint + "/" + subscriptionId.String())
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
	var subscriptions models.SubscriptionDetail
	err = json.Unmarshal(response, &subscriptions)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &subscriptions, nil
}

func (pc tmsClient) GetSubscriptionPolicies(serviceId, subscriptionId uuid.UUID) (*models.SubscriptionPolicies, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.ServiceApiEndpoint + "/" +
		serviceId.String() + constants.SubscriptionApiEndpoint + "/" + subscriptionId.String() + constants.PolicyApiEndpoint)
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
	var subscriptions models.SubscriptionPolicies
	err = json.Unmarshal(response, &subscriptions)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &subscriptions, nil
}

func (pc tmsClient) GetSubscriptionTagValues(serviceId, subscriptionId uuid.UUID) (*models.SubscriptionTagsValues, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.ServiceApiEndpoint + "/" +
		serviceId.String() + constants.SubscriptionApiEndpoint + "/" + subscriptionId.String() + constants.TagsValuesEndpoint)
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
	var subscriptions models.SubscriptionTagsValues
	err = json.Unmarshal(response, &subscriptions)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &subscriptions, nil
}

func (pc tmsClient) DeleteSubscription(serviceId, subscriptionId uuid.UUID) error {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.ServiceApiEndpoint + "/" +
		serviceId.String() + constants.SubscriptionApiEndpoint + "/" + subscriptionId.String())
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

func (pc tmsClient) CreateUser(user *models.CreateTenantUser) (*models.User, error) {
	reqBytes, err := json.Marshal(user)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling request")
	}

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.UserApiEndpoint)
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
	var createUserRes models.User
	err = json.Unmarshal(response, &createUserRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &createUserRes, nil
}

func (pc tmsClient) UpdateUser(user *models.UpdateUser) (*models.User, error) {
	reqBytes, err := json.Marshal(user)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling request")
	}

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.UserApiEndpoint + "/" + user.Id.String())
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
	var createUserRes models.User
	err = json.Unmarshal(response, &createUserRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &createUserRes, nil
}

func (pc tmsClient) UpdateTenantUserRole(request *models.UpdateTenantUserRoles) (*models.User, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling request")
	}

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.UserApiEndpoint +
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
	var updateUserRes models.User
	err = json.Unmarshal(response, &updateUserRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &updateUserRes, nil
}

func (pc tmsClient) GetUsers() ([]models.User, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.UserApiEndpoint)
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
	var searchUserRes []models.User
	err = json.Unmarshal(response, &searchUserRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return searchUserRes, nil
}

func (pc tmsClient) RetrieveUser(id uuid.UUID) (*models.User, error) {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.UserApiEndpoint + "/" + id.String())
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
	var retrieveUserRes *models.User
	err = json.Unmarshal(response, &retrieveUserRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return retrieveUserRes, nil
}

func (pc tmsClient) DeleteUser(userId uuid.UUID) error {
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.UserApiEndpoint + "/" + userId.String())
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

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.ServiceApiEndpoint)
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

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.ServiceApiEndpoint + "/" + request.Id.String())
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

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.ServiceApiEndpoint + "/" + serviceId.String())
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
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.ServiceApiEndpoint)
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
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.ServiceApiEndpoint + "/" + id.String())
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
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.TagApiEndpoint)
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
	reqURL, err := url.Parse(pc.BaseURL.String() + constants.TenantApiEndpoint + constants.TagApiEndpoint)
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

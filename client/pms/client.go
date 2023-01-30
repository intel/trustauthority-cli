/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package pms

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	client "intel/amber/tac/v1/client"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/models"
	"net/http"
	"net/url"
)

type PmsClient interface {
	CreatePolicy(policyRequest *models.PolicyRequest) (*models.PolicyResponse, error)
	DeletePolicy(policyID uuid.UUID) error
	GetPolicy(policyID uuid.UUID) (*models.PolicyResponse, error)
	UpdatePolicy(request *models.PolicyUpdateRequest) (*models.PolicyResponse, error)
	SearchPolicy() ([]models.PolicyResponse, error)
}

//Client Details for PMS client
type pmsClient struct {
	Client   *http.Client
	BaseURL  *url.URL
	TenantId uuid.UUID
	ApiKey   string
}

func NewPmsClient(client *http.Client, pmsURL *url.URL, tenantId uuid.UUID, apiKey string) PmsClient {
	return &pmsClient{
		Client:   client,
		BaseURL:  pmsURL,
		TenantId: tenantId,
		ApiKey:   apiKey,
	}
}

func (pc pmsClient) CreatePolicy(request *models.PolicyRequest) (*models.PolicyResponse, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling request")
	}

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.PolicyApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL: %s", pc.BaseURL.String())
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPost, reqURL.String(), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response
	var policyRes models.PolicyResponse
	err = json.Unmarshal(response, &policyRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &policyRes, nil
}

func (pc pmsClient) DeletePolicy(policyID uuid.UUID) error {

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.PolicyApiEndpoint + "/" + policyID.String())
	if err != nil {
		return errors.Wrapf(err, "Invalid URL: %s", pc.BaseURL.String())
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

	return nil
}

func (pc pmsClient) GetPolicy(policyID uuid.UUID) (*models.PolicyResponse, error) {

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.PolicyApiEndpoint + "/" + policyID.String())
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL: %s", pc.BaseURL.String())
	}

	log.Debugf("PMS Request URL: %s", reqURL)

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

	// Parse response
	var policyRes models.PolicyResponse
	err = json.Unmarshal(response, &policyRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling policy response")
	}
	return &policyRes, nil
}

func (pc pmsClient) SearchPolicy() ([]models.PolicyResponse, error) {

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.PolicyApiEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL: %s", pc.BaseURL.String())
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

	// Parse response
	var policyRes []models.PolicyResponse
	err = json.Unmarshal(response, &policyRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return policyRes, nil
}

func (pc pmsClient) UpdatePolicy(request *models.PolicyUpdateRequest) (*models.PolicyResponse, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, " Error marshalling policy update request")
	}

	reqURL, err := url.Parse(pc.BaseURL.String() + constants.PolicyApiEndpoint + "/" + request.PolicyId.String())
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL: %s", pc.BaseURL.String())
	}

	log.Debugf("PMS Request URL: %s", reqURL)

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPut, reqURL.String(), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errors.Wrap(err, " Error forming request")
	}
	req.Header.Add(constants.HTTPHeaderKeyAccept, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyContentType, constants.HTTPMediaTypeJson)
	req.Header.Add(constants.HTTPHeaderKeyApiKey, pc.ApiKey)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response
	var policyRes models.PolicyResponse
	err = json.Unmarshal(response, &policyRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &policyRes, nil
}

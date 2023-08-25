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
	"intel/tac/v1/client"
	"intel/tac/v1/constants"
	models2 "intel/tac/v1/internal/models"
	"intel/tac/v1/models"
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

// Client Details for PMS client
type pmsClient struct {
	Client  *http.Client
	BaseURL *url.URL
	ApiKey  string
}

func NewPmsClient(client *http.Client, pmsURL *url.URL, apiKey string) PmsClient {
	return &pmsClient{
		Client:  client,
		BaseURL: pmsURL,
		ApiKey:  apiKey,
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
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response
	var policyRes models.PolicyResponse
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&policyRes)

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
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

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
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response
	var policyRes models.PolicyResponse
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&policyRes)
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
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response
	var policyRes []models.PolicyResponse
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&policyRes)
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
	req.Header.Add(constants.HTTPHeaderKeyRequestId, models2.RespHeaderFields.RequestId)

	response, err := client.SendRequest(pc.Client, req)
	if err != nil {
		return nil, errors.Wrap(err, "Error in response body")
	}

	// Parse response
	var policyRes models.PolicyResponse
	dec := json.NewDecoder(bytes.NewReader(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&policyRes)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling response")
	}
	return &policyRes, nil
}

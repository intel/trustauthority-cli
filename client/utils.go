/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package client

import (
	"context"
	"fmt"
	rClient "github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"intel/tac/v1/constants"
	"intel/tac/v1/internal/models"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	retryableStatusCode = map[int]bool{
		500: true,
		503: true,
		504: true,
	}
)

func SendRequest(client *http.Client, req *http.Request) ([]byte, error) {
	var resp *http.Response
	var err error

	// set the request Id to the provided in case there is an error while sending and receiving request
	models.RespHeaderFields.RequestId = req.Header.Get(constants.HTTPHeaderKeyRequestId)

	var retryClient = rClient.NewClient()
	retryClient.HTTPClient = client
	retryClient.RetryWaitMin = constants.DefaultRetryWaitMin * time.Second
	retryClient.RetryWaitMax = constants.DefaultRetryWaitMax * time.Second
	retryClient.RetryMax = constants.DefaultRetryCount
	retryClient.CheckRetry = retryPolicy
	retryClient.Logger = log.StandardLogger()

	if resp, err = retryClient.StandardClient().Do(req); err != nil {
		return nil, err
	}

	if resp != nil {
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				log.WithError(err).WithField(constants.HTTPHeaderKeyRequestId, req.Header.Get(constants.HTTPHeaderKeyRequestId)).Errorf("Failed to close response body")
			}
		}()
	}

	//Get the request and trace ID from response header
	models.RespHeaderFields.RequestId = resp.Header.Get(constants.HTTPHeaderKeyRequestId)
	models.RespHeaderFields.TraceId = resp.Header.Get(constants.HTTPHeaderKeyTraceId)

	//create byte array of HTTP response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response")
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, errors.Errorf("The call to %q returned %q. Error: %s", req.URL, resp.Status, body)
	}
	return body, nil
}

func retryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// Do not retry on context.Canceled
	if ctx.Err() != nil {
		// If connection was closed due to client timeout retry again
		if ctx.Err() == context.DeadlineExceeded {
			return true, ctx.Err()
		}
		return false, nil
	}

	//Retry if the request did not reach the API gateway and the error is Service Unavailable
	if err != nil {
		if v, ok := err.(*url.Error); ok {
			if strings.ToLower(v.Error()) == constants.ServiceUnavailableError {
				return true, v
			}
		}
		return false, nil
	}

	// Check the response code. We retry on 500, 503 and 504 responses to allow
	// the server time to recover, as these are typically not permanent
	// errors and may relate to outages on the server side.
	if ok := retryableStatusCode[resp.StatusCode]; ok {
		return true, fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}
	return false, nil
}

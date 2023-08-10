/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package client

import (
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/internal/models"
	"io"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func SendRequest(client *http.Client, req *http.Request) ([]byte, error) {
	var resp *http.Response
	var err error

	// set the request Id to the provided in case there is an error while sending and receiving request
	models.RespHeaderFields.RequestId = req.Header.Get(constants.HTTPHeaderKeyRequestId)
	if resp, err = client.Do(req); err != nil {
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

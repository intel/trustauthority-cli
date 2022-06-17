/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package client

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func SendRequest(client *http.Client, req *http.Request) ([]byte, error) {
	var resp *http.Response
	var err error

	if resp, err = client.Do(req); err != nil {
		return nil, err
	}

	defer func() {
		derr := resp.Body.Close()
		if derr != nil {
			log.WithError(derr).Error("Error closing response body")
		}
	}()

	if resp != nil {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.WithError(err).Errorf("Failed to close response body")
			}
		}()
	}

	//create byte array of HTTP response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response")
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, errors.Errorf("The call to %q return %q", req.URL, resp.Status)
	}

	return body, nil
}

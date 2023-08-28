/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"intel/tac/v1/client/tms"
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	"intel/tac/v1/utils"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// getApiClientsCmd represents the getApiClients command
var getApiClientsCmd = &cobra.Command{
	Use:   constants.ApiClientCmd,
	Short: "Get apiClient(s) under a particular tenant",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("list apiClients called")
		response, err := getApiClients(cmd)
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Println("ApiClients: \n\n", response)
		return nil
	},
}

func init() {
	listCmd.AddCommand(getApiClientsCmd)

	getApiClientsCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Trust Authority service for which the apiClient needs to be created")
	getApiClientsCmd.Flags().StringP(constants.ApiClientIdParamName, "c", "", "Id of the apiClient which needs to be fetched (optional)")
	getApiClientsCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
	getApiClientsCmd.MarkFlagRequired(constants.ServiceIdParamName)
}

func getApiClients(cmd *cobra.Command) (string, error) {
	configValues, err := config.LoadConfiguration()
	if err != nil {
		return "", err
	}
	client := &http.Client{
		Timeout: time.Duration(configValues.HTTPClientTimeout) * time.Second,
	}

	tmsUrl, err := url.Parse(configValues.TrustAuthorityBaseUrl + constants.TmsBaseUrl)
	if err != nil {
		return "", err
	}

	if err = setRequestId(cmd); err != nil {
		return "", err
	}

	serviceIdString, err := cmd.Flags().GetString(constants.ServiceIdParamName)
	if err != nil {
		return "", err
	}

	serviceId, err := uuid.Parse(serviceIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid service id provided")
	}

	apiClientIdString, err := cmd.Flags().GetString(constants.ApiClientIdParamName)
	if err != nil {
		return "", err
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, apiKey)

	var responseBytes []byte
	if apiClientIdString == "" {
		fmt.Println("API client ID is not set, fetching all API clients ...")
		response, err := tmsClient.GetApiClient(serviceId)
		if err != nil {
			return "", err
		}

		responseBytes, err = json.MarshalIndent(response, "", "  ")
		if err != nil {
			return "", err
		}
	} else {
		apiClientId, err := uuid.Parse(apiClientIdString)
		if err != nil {
			return "", errors.Wrap(err, "Invalid apiClient id provided")
		}

		response, err := tmsClient.RetrieveApiClient(serviceId, apiClientId)
		if err != nil {
			return "", err
		}

		responseBytes, err = json.MarshalIndent(response, "", "  ")
		if err != nil {
			return "", err
		}
	}

	return string(responseBytes), nil
}

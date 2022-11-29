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
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
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
		if err != nil {
			return err
		}
		fmt.Println("ApiClients: \n\n", response)
		return nil
	},
}

func init() {
	listCmd.AddCommand(getApiClientsCmd)

	getApiClientsCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	getApiClientsCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Amber service for which the apiClient needs to be created")
	getApiClientsCmd.Flags().StringP(constants.ApiClientIdParamName, "c", "", "Id of the apiClient which needs to be fetched (optional)")
	getApiClientsCmd.MarkFlagRequired(constants.ApiKeyParamName)
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

	tmsUrl, err := url.Parse(configValues.AmberBaseUrl + constants.TmsBaseUrl)
	if err != nil {
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

	tmsClient := tms.NewTmsClient(client, tmsUrl, uuid.Nil, apiKey)

	var responseBytes []byte
	if apiClientIdString == "" {
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

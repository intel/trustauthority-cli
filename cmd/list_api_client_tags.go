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

var getApiClientTagsValuesCmd = &cobra.Command{
	Use:   constants.TagCmd,
	Short: "Used to get the list of tags and their values IDs linked to a apiClient",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("list apiClient tag called")
		response, err := getApiClientTagsAndValues(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Services: \n\n", response)
		return nil
	},
}

func init() {
	getApiClientsCmd.AddCommand(getApiClientTagsValuesCmd)

	getApiClientTagsValuesCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	getApiClientTagsValuesCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Amber service for which the apiClient policies are to be fetched")
	getApiClientTagsValuesCmd.Flags().StringP(constants.ApiClientIdParamName, "c", "", "Id of the apiClient for which the policies are to be fetched")
	getApiClientTagsValuesCmd.MarkFlagRequired(constants.ApiKeyParamName)
	getApiClientTagsValuesCmd.MarkFlagRequired(constants.ServiceIdParamName)
	getApiClientTagsValuesCmd.MarkFlagRequired(constants.ApiClientIdParamName)
}

func getApiClientTagsAndValues(cmd *cobra.Command) (string, error) {
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
	apiClientId, err := uuid.Parse(apiClientIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid apiClient id provided")
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, uuid.Nil, apiKey)
	response, err := tmsClient.GetApiClientTagValues(serviceId, apiClientId)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

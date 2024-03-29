/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 *
 *
 */

package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"intel/tac/v1/client/tms"
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	"intel/tac/v1/utils"
	"net/http"
	"net/url"
	"time"
)

var deleteApiClientCmd = &cobra.Command{
	Use:   constants.ApiClientCmd,
	Short: "Delete an api client whose ID has been provided",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("delete apiClient called")
		serviceId, err := deleteApiClient(cmd)
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Printf("Deleted api client with Id: %s \n\n", serviceId)
		fmt.Println("\nNOTE: There may be a delay of up to two (2) minutes for the changes to the attestation API key to take effect.")
		fmt.Print("\n")
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteApiClientCmd)

	deleteApiClientCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Trust Authority service for which the api client needs to be created")
	deleteApiClientCmd.Flags().StringP(constants.ApiClientIdParamName, "c", "", "Id of the api client which needs to be fetched (optional)")
	deleteApiClientCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
	deleteApiClientCmd.MarkFlagRequired(constants.ServiceIdParamName)
	deleteApiClientCmd.MarkFlagRequired(constants.ApiClientIdParamName)
}

func deleteApiClient(cmd *cobra.Command) (string, error) {
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

	apiClientId, err := uuid.Parse(apiClientIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid api client id provided")
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, apiKey)

	err = tmsClient.DeleteApiClient(serviceId, apiClientId)
	if err != nil {
		return "", err
	}

	return apiClientIdString, nil
}

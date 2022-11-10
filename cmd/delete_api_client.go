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
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
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
		if err != nil {
			return err
		}
		fmt.Printf("Deleted api client with Id: %s \n\n", serviceId)
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteApiClientCmd)

	deleteApiClientCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	deleteApiClientCmd.Flags().StringP(constants.TenantIdParamName, "t", "", "Id of the tenant for whom the api client needs to be created")
	deleteApiClientCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Amber service for which the api client needs to be created")
	deleteApiClientCmd.Flags().StringP(constants.ApiClientIdParamName, "d", "", "Id of the api client which needs to be fetched (optional)")
	deleteApiClientCmd.MarkFlagRequired(constants.ApiKeyParamName)
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

	tmsUrl, err := url.Parse(configValues.AmberBaseUrl + constants.TmsBaseUrl)
	if err != nil {
		return "", err
	}

	tenantIdString, err := cmd.Flags().GetString(constants.TenantIdParamName)
	if err != nil {
		return "", err
	}

	if tenantIdString == "" {
		tenantIdString = configValues.TenantId
	}

	tenantId, err := uuid.Parse(tenantIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid tenant id provided")
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

	tmsClient := tms.NewTmsClient(client, tmsUrl, tenantId, apiKey)

	apiClientId, err := uuid.Parse(apiClientIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid api client id provided")
	}

	err = tmsClient.DeleteApiClient(serviceId, apiClientId)
	if err != nil {
		return "", err
	}

	return apiClientIdString, nil
}

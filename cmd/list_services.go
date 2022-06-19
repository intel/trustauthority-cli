/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// getServicesCmd represents the getServices command
var getServicesCmd = &cobra.Command{
	Use:   "service",
	Short: "Used to get the list of services under a tenant",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("get services called")
		response, err := getServices(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Services: \n\n", response)
		return nil
	},
}

func init() {
	listCmd.AddCommand(getServicesCmd)

	getServicesCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	getServicesCmd.Flags().StringP(constants.TenantIdParamName, "t", "", "Id of the tenant for whom the subscription needs to be created")
	getServicesCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Amber service for which the subscription needs to be created")
	getServicesCmd.MarkFlagRequired(constants.ApiKeyParamName)
}

func getServices(cmd *cobra.Command) (string, error) {
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
		return "", err
	}

	serviceIdString, err := cmd.Flags().GetString(constants.ServiceIdParamName)
	if err != nil {
		return "", err
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, tenantId, apiKey)

	var responseBytes []byte
	if serviceIdString == "" {
		response, err := tmsClient.GetServices()
		if err != nil {
			return "", err
		}

		responseBytes, err = json.MarshalIndent(response, "", "  ")
		if err != nil {
			return "", err
		}
	} else {
		serviceId, err := uuid.Parse(serviceIdString)
		if err != nil {
			return "", err
		}

		response, err := tmsClient.RetrieveService(serviceId)
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

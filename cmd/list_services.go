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

// getServicesCmd represents the getServices command
var getServicesCmd = &cobra.Command{
	Use:   constants.ServiceCmd,
	Short: "Used to get the list of services under a tenant",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("list services called")
		response, err := getServices(cmd)
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Println("Services: \n\n", response)
		return nil
	},
}

func init() {
	listCmd.AddCommand(getServicesCmd)

	getServicesCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Trust Authority service which needs to be fetched")
	getServicesCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
}

func getServices(cmd *cobra.Command) (string, error) {
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

	tmsClient := tms.NewTmsClient(client, tmsUrl, apiKey)

	var responseBytes []byte
	if serviceIdString == "" {
		fmt.Println("Service ID was not provided, listing all services....")
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
			return "", errors.Wrap(err, "Invalid service id provided")
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

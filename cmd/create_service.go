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
	"intel/amber/tac/v1/models"
	"intel/amber/tac/v1/validation"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// createServiceCmd represents the createService command
var createServiceCmd = &cobra.Command{
	Use:   constants.ServiceCmd,
	Short: "Create a new service",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("create service called")
		response, err := createService(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Service: \n\n", response)
		return nil
	},
}

func init() {
	createCmd.AddCommand(createServiceCmd)

	createServiceCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	createServiceCmd.Flags().StringP(constants.TenantIdParamName, "t", "", "Id of the tenant for whom the subscription needs to be created")
	createServiceCmd.Flags().StringP(constants.ServiceOfferIdParamName, "r", "", "Id of the Amber service offer for which the service needs to be created")
	createServiceCmd.Flags().StringP(constants.ServiceNameParamName, "n", "", "Description of the service")
	createServiceCmd.MarkFlagRequired(constants.ApiKeyParamName)
	createServiceCmd.MarkFlagRequired(constants.ServiceOfferIdParamName)
}

func createService(cmd *cobra.Command) (string, error) {
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

	serviceOfferIdString, err := cmd.Flags().GetString(constants.ServiceOfferIdParamName)
	if err != nil {
		return "", err
	}

	serviceOfferId, err := uuid.Parse(serviceOfferIdString)
	if err != nil {
		return "", err
	}

	serviceDescription, err := cmd.Flags().GetString(constants.ServiceNameParamName)
	if err != nil {
		return "", err
	}

	if err = validation.ValidateStrings([]string{serviceDescription}); err != nil {
		return "", errors.Wrap(err, "Invalid service description provided")
	}
	var serviceInfo = models.CreateService{
		ServiceOfferId: serviceOfferId,
		Description:    serviceDescription,
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, tenantId, apiKey)
	response, err := tmsClient.CreateService(&serviceInfo)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

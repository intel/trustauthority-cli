/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 *
 *
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/models"
	"net/http"
	"net/url"
	"time"
)

var updateServiceCmd = &cobra.Command{
	Use:   constants.ServiceCmd,
	Short: "Update a service whose ID has been provided",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("update service called")
		response, err := updateService(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Updated service: \n\n", response)
		return nil
	},
}

func init() {
	updateCmd.AddCommand(updateServiceCmd)

	updateServiceCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	updateServiceCmd.Flags().StringP(constants.ServiceIdParamName, "s", "", "Id of the service to be updated")
	updateServiceCmd.Flags().StringP(constants.ServiceNameParamName, "n", "", "Name of the service ")
	updateServiceCmd.MarkFlagRequired(constants.ApiKeyParamName)
	updateServiceCmd.MarkFlagRequired(constants.ServiceIdParamName)
	updateServiceCmd.MarkFlagRequired(constants.ServiceNameParamName)
}

func updateService(cmd *cobra.Command) (string, error) {
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
		return "", errors.Wrap(err, "Invalid service Id provided")
	}

	serviceName, err := cmd.Flags().GetString(constants.ServiceNameParamName)
	if err != nil {
		return "", err
	}

	var serviceUpdateReq = &models.UpdateService{
		Id:   serviceId,
		Name: serviceName,
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, uuid.Nil, apiKey)
	response, err := tmsClient.UpdateService(serviceUpdateReq)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

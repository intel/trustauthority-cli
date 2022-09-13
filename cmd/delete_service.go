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

var deleteServiceCmd = &cobra.Command{
	Use:   constants.ServiceCmd,
	Short: "Delete a service whose ID has been provided",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("delete service called")
		serviceId, err := deleteService(cmd)
		if err != nil {
			return err
		}
		fmt.Printf("Deleted service with Id: %s \n\n", serviceId)
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteServiceCmd)

	deleteServiceCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	deleteServiceCmd.Flags().StringP(constants.ServiceIdParamName, "s", "", "Id of the service to be deleted")
	deleteServiceCmd.MarkFlagRequired(constants.ApiKeyParamName)
	deleteServiceCmd.MarkFlagRequired(constants.ServiceIdParamName)
}

func deleteService(cmd *cobra.Command) (string, error) {
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

	tmsClient := tms.NewTmsClient(client, tmsUrl, uuid.Nil, apiKey)
	err = tmsClient.DeleteService(serviceId)
	if err != nil {
		return "", err
	}

	return serviceIdString, nil
}

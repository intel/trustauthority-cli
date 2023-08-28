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

var getApiClientPoliciesCmd = &cobra.Command{
	Use:   constants.PolicyCmd,
	Short: "Used to get the list of policy IDs linked to a apiClient",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("list apiClient policy called")
		response, err := getApiClientPolicies(cmd)
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Println("Policy IDs: \n\n", response)
		return nil
	},
}

func init() {
	getApiClientsCmd.AddCommand(getApiClientPoliciesCmd)

	getApiClientPoliciesCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Trust Authority service for which the apiClient policies are to be fetched")
	getApiClientPoliciesCmd.Flags().StringP(constants.ApiClientIdParamName, "c", "", "Id of the apiClient for which the policies are to be fetched")
	getApiClientPoliciesCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
	getApiClientPoliciesCmd.MarkFlagRequired(constants.ServiceIdParamName)
	getApiClientPoliciesCmd.MarkFlagRequired(constants.ApiClientIdParamName)
}

func getApiClientPolicies(cmd *cobra.Command) (string, error) {
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
		return "", errors.Wrap(err, "Invalid apiClient id provided")
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, apiKey)
	response, err := tmsClient.GetApiClientPolicies(serviceId, apiClientId)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

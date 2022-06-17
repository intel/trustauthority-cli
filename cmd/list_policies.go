/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"intel/amber/tac/v1/client/pms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// getPoliciesCmd represents the getPolicies command
var getPoliciesCmd = &cobra.Command{
	Use:   "policy",
	Short: "Get list of policies or specific policy user a tenant",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("get policies called")
		response, err := getPolicies(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Policies: \n\n", response)
		return nil
	},
}

func init() {
	listCmd.AddCommand(getPoliciesCmd)

	getPoliciesCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	getPoliciesCmd.Flags().StringP(constants.PolicyIdParamName, "p", "", "Path of the file containing the policy to be uploaded")
	getPoliciesCmd.Flags().StringP(constants.TenantIdParamName, "t", "", "Id of the tenant for whom the policies need to be fetched")
	getPoliciesCmd.MarkFlagRequired(constants.ApiKeyParamName)
	getPoliciesCmd.MarkFlagRequired(constants.TenantIdParamName)
}

func getPolicies(cmd *cobra.Command) (string, error) {

	configValues, err := config.LoadConfiguration()
	if err != nil {
		return "", err
	}
	client := &http.Client{
		Timeout: time.Duration(configValues.HTTPClientTimeout) * time.Second,
	}

	pmsUrl, err := url.Parse(configValues.AmberBaseUrl + constants.PmsBaseUrl)
	if err != nil {
		return "", err
	}

	tenantIdString, err := cmd.Flags().GetString(constants.TenantIdParamName)
	if err != nil {
		return "", err
	}

	tenantId, err := uuid.Parse(tenantIdString)
	if err != nil {
		return "", err
	}

	policyIdString, err := cmd.Flags().GetString(constants.PolicyIdParamName)
	if err != nil {
		return "", err
	}

	pmsClient := pms.NewPmsClient(client, pmsUrl, tenantId, apiKey)

	var responseBytes []byte
	if policyIdString == "" {
		response, err := pmsClient.SearchPolicy()
		responseBytes, err = json.MarshalIndent(response, "", "  ")
		if err != nil {
			return "", err
		}
	} else {
		policyId, err := uuid.Parse(policyIdString)
		if err != nil {
			return "", err
		}
		response, err := pmsClient.GetPolicy(policyId)
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

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
	"intel/amber/tac/v1/models"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// createPolicyCmd represents the createPolicy command
var createPolicyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Create a new policy",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("create policy called")
		response, err := createPolicy(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Policy: \n\n", response)
		return nil
	},
}

func init() {
	createCmd.AddCommand(createPolicyCmd)

	createPolicyCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	createPolicyCmd.Flags().StringP(constants.PolicyFileParamName, "f", "", "Path of the file containing the policy to be uploaded")
	createPolicyCmd.MarkFlagRequired(constants.ApiKeyParamName)
	createPolicyCmd.MarkFlagRequired(constants.PolicyFileParamName)
}

func createPolicy(cmd *cobra.Command) (string, error) {
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

	policyFilePath, err := cmd.Flags().GetString(constants.PolicyFileParamName)
	if err != nil {
		return "", err
	}

	policyBytes, err := os.ReadFile(policyFilePath)
	if err != nil {
		return "", err
	}

	var policyCreateReq models.PolicyRequest
	err = json.Unmarshal(policyBytes, &policyCreateReq)
	if err != nil {
		return "", err
	}

	pmsClient := pms.NewPmsClient(client, pmsUrl, uuid.Nil, apiKey)
	response, err := pmsClient.CreatePolicy(&policyCreateReq)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

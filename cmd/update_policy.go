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

// updatePolicyCmd represents the updatePolicy command
var updatePolicyCmd = &cobra.Command{
	Use:   constants.PolicyCmd,
	Short: "Update a policy whose ID has been provided",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("update Policy called")
		response, err := updatePolicy(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Updated policy: \n\n", response)
		return nil
	},
}

func init() {
	updateCmd.AddCommand(updatePolicyCmd)

	updatePolicyCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	updatePolicyCmd.Flags().StringP(constants.PolicyFileParamName, "f", "", "Path of the file containing the policy to be uploaded")
	updatePolicyCmd.MarkFlagRequired(constants.ApiKeyParamName)
	updatePolicyCmd.MarkFlagRequired(constants.PolicyFileParamName)
}

func updatePolicy(cmd *cobra.Command) (string, error) {
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

	var policyUpdateReq models.PolicyRequest
	err = json.Unmarshal(policyBytes, &policyUpdateReq)
	if err != nil {
		return "", err
	}

	if policyUpdateReq.PolicyId == uuid.Nil {
		return "", errors.Errorf("Please add the policy_id to the JSON request in policy file %s", policyFilePath)
	}

	pmsClient := pms.NewPmsClient(client, pmsUrl, uuid.Nil, apiKey)
	response, err := pmsClient.UpdatePolicy(&policyUpdateReq)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

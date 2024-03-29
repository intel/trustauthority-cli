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
	"intel/tac/v1/client/pms"
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	"intel/tac/v1/models"
	"intel/tac/v1/utils"
	"intel/tac/v1/validation"
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
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Println("Updated policy: \n\n", response)
		return nil
	},
}

func init() {
	updateCmd.AddCommand(updatePolicyCmd)

	updatePolicyCmd.Flags().StringP(constants.PolicyIdParamName, "i", "", "Id of the policy to be updated")
	updatePolicyCmd.Flags().StringP(constants.PolicyNameParamName, "n", "", "Name of the policy to be updated")
	updatePolicyCmd.Flags().StringP(constants.PolicyFileParamName, "f", "", "Path of the file containing the rego policy to be uploaded. The file size should be <= 10 KB")
	updatePolicyCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
	updatePolicyCmd.MarkFlagRequired(constants.PolicyIdParamName)
}

func updatePolicy(cmd *cobra.Command) (string, error) {
	configValues, err := config.LoadConfiguration()
	if err != nil {
		return "", err
	}
	client := &http.Client{
		Timeout: time.Duration(configValues.HTTPClientTimeout) * time.Second,
	}

	pmsUrl, err := url.Parse(configValues.TrustAuthorityBaseUrl + constants.PmsBaseUrl)
	if err != nil {
		return "", err
	}

	if err = setRequestId(cmd); err != nil {
		return "", err
	}

	policyIdString, err := cmd.Flags().GetString(constants.PolicyIdParamName)
	if err != nil {
		return "", err
	}

	policyId, err := uuid.Parse(policyIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid policy Id provided, should be in UUID format")
	}

	var policyUpdateReq = models.PolicyUpdateRequest{PolicyId: policyId}

	policyName, err := cmd.Flags().GetString(constants.PolicyNameParamName)
	if err != nil {
		return "", err
	}

	if policyName != "" {
		if err = validation.ValidatePolicyName(policyName); err != nil {
			return "", err
		}
		policyUpdateReq.PolicyName = policyName
	}

	policyFilePath, err := cmd.Flags().GetString(constants.PolicyFileParamName)
	if err != nil {
		return "", err
	}
	// policy file is not mandatory, skipping policy read if file path is empty
	if policyFilePath != "" {
		path, err := validation.ValidatePath(policyFilePath)
		if err != nil {
			return "", err
		}
		policyBytes, err := os.ReadFile(path)
		if err != nil {
			return "", errors.Wrap(err, "Error reading policy file")
		}

		err = validation.ValidateSize(policyFilePath)
		if err != nil {
			return "", err
		}

		if string(policyBytes) != "" {
			policyUpdateReq.Policy = string(policyBytes)
		}
	}

	pmsClient := pms.NewPmsClient(client, pmsUrl, apiKey)
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

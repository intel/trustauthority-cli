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

// createPolicyCmd represents the createPolicy command
var createPolicyCmd = &cobra.Command{
	Use:   constants.PolicyCmd,
	Short: "Create a new policy",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("create policy called")
		response, err := createPolicy(cmd)
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Println("Policy: \n\n", response)
		return nil
	},
}

func init() {
	createCmd.AddCommand(createPolicyCmd)

	createPolicyCmd.Flags().StringP(constants.PolicyNameParamName, "n", "", "Name of the policy to be uploaded")
	createPolicyCmd.Flags().StringP(constants.PolicyTypeParamName, "t", "", "Type of the policy to be uploaded, example \"Appraisal policy\".")
	createPolicyCmd.Flags().StringP(constants.ServiceOfferIdParamName, "r", "", "Service offer id for which the policy needs to be uploaded")
	createPolicyCmd.Flags().StringP(constants.AttestationTypeParamName, "a", "", "Attestation type of policy to be uploaded, example \"SGX Attestation\".")
	createPolicyCmd.Flags().StringP(constants.PolicyFileParamName, "f", "", "Path of the file containing the rego policy to be uploaded. The file size should be <= 10 KB")
	createPolicyCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
	createPolicyCmd.MarkFlagRequired(constants.PolicyNameParamName)
	createPolicyCmd.MarkFlagRequired(constants.ServiceOfferIdParamName)
	createPolicyCmd.MarkFlagRequired(constants.AttestationTypeParamName)
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

	pmsUrl, err := url.Parse(configValues.TrustAuthorityBaseUrl + constants.PmsBaseUrl)
	if err != nil {
		return "", err
	}

	if err = setRequestId(cmd); err != nil {
		return "", err
	}

	policyName, err := cmd.Flags().GetString(constants.PolicyNameParamName)
	if err != nil {
		return "", err
	}
	if err = validation.ValidatePolicyName(policyName); err != nil {
		return "", err
	}

	policyType, err := cmd.Flags().GetString(constants.PolicyTypeParamName)
	if err != nil {
		return "", err
	}

	soIdString, err := cmd.Flags().GetString(constants.ServiceOfferIdParamName)
	if err != nil {
		return "", err
	}

	soId, err := uuid.Parse(soIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid service offer Id provided, should be in UUID format")
	}

	attestationType, err := cmd.Flags().GetString(constants.AttestationTypeParamName)
	if err != nil {
		return "", err
	}

	policyFilePath, err := cmd.Flags().GetString(constants.PolicyFileParamName)
	if err != nil {
		return "", err
	}
	if policyFilePath == "" {
		return "", errors.New("Policy file path cannot be empty")
	}

	path, err := validation.ValidatePath(policyFilePath)
	if err != nil {
		return "", err
	}

	err = validation.ValidateSize(policyFilePath)
	if err != nil {
		return "", err
	}
	policyBytes, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "Error reading policy file")
	}

	var policyCreateReq = models.PolicyRequest{models.CommonPolicy{
		Policy:          string(policyBytes),
		PolicyName:      policyName,
		PolicyType:      policyType,
		ServiceOfferId:  soId,
		AttestationType: attestationType,
	}}

	pmsClient := pms.NewPmsClient(client, pmsUrl, apiKey)
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

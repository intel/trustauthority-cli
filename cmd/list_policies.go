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
	"intel/tac/v1/utils"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// getPoliciesCmd represents the getPolicies command
var getPoliciesCmd = &cobra.Command{
	Use:   constants.PolicyCmd,
	Short: "Get list of policies or specific policy user a tenant",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("list policies called")
		response, err := getPolicies(cmd)
		if err != nil {
			return err
		}
		utils.PrintRequestAndTraceId()
		fmt.Println("Policies: \n\n", response)
		return nil
	},
}

func init() {
	listCmd.AddCommand(getPoliciesCmd)

	getPoliciesCmd.Flags().StringP(constants.PolicyIdParamName, "p", "", "Path of the file containing the policy to be uploaded")
	getPoliciesCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
}

func getPolicies(cmd *cobra.Command) (string, error) {

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

	pmsClient := pms.NewPmsClient(client, pmsUrl, apiKey)

	var responseBytes []byte
	if policyIdString == "" {
		response, err := pmsClient.SearchPolicy()
		if err != nil {
			return "", err
		}
		responseBytes, err = json.MarshalIndent(response, "", "  ")
		if err != nil {
			return "", err
		}
	} else {
		policyId, err := uuid.Parse(policyIdString)
		if err != nil {
			return "", errors.Wrap(err, "Invalid policy id provided")
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

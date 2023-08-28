/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
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

// deletePolicyCmd represents the deletePolicy command
var deletePolicyCmd = &cobra.Command{
	Use:   constants.PolicyCmd,
	Short: "Deletes a policy",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("delete policy called")
		policyId, err := deletePolicy(cmd)
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Printf("\nPolicy %s deleted", policyId)
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deletePolicyCmd)

	deletePolicyCmd.Flags().StringP(constants.PolicyIdParamName, "p", "", "Id of the policy to be deleted")
	deletePolicyCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
	deletePolicyCmd.MarkFlagRequired(constants.PolicyIdParamName)
}

func deletePolicy(cmd *cobra.Command) (string, error) {
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
		return "", errors.Wrap(err, "Invalid policy id provided")
	}

	pmsClient := pms.NewPmsClient(client, pmsUrl, apiKey)

	err = pmsClient.DeletePolicy(policyId)
	if err != nil {
		return "", err
	}

	return policyIdString, nil
}

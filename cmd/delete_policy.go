/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"intel/amber/tac/v1/client/pms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// deletePolicyCmd represents the deletePolicy command
var deletePolicyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Deletes a policy",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("delete policy called")
		policyId, err := deletePolicy(cmd)
		if err != nil {
			return err
		}
		fmt.Printf("\nPolicy %s deleted", policyId)
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deletePolicyCmd)

	deletePolicyCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	deletePolicyCmd.Flags().StringP(constants.PolicyIdParamName, "p", "", "Id of the policy to be deleted")
	deletePolicyCmd.Flags().StringP(constants.TenantIdParamName, "t", "", "Id of the tenant for whom the subscription needs to be created")
	deletePolicyCmd.MarkFlagRequired(constants.ApiKeyParamName)
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

	pmsUrl, err := url.Parse(configValues.AmberBaseUrl + constants.PmsBaseUrl)
	if err != nil {
		return "", err
	}

	policyIdString, err := cmd.Flags().GetString(constants.PolicyIdParamName)
	if err != nil {
		return "", err
	}

	policyId, err := uuid.Parse(policyIdString)
	if err != nil {
		return "", err
	}

	tenantIdString, err := cmd.Flags().GetString(constants.TenantIdParamName)
	if err != nil {
		return "", err
	}

	if tenantIdString == "" {
		tenantIdString = configValues.TenantId
	}

	tenantId, err := uuid.Parse(tenantIdString)
	if err != nil {
		return "", err
	}

	pmsClient := pms.NewPmsClient(client, pmsUrl, tenantId, apiKey)

	err = pmsClient.DeletePolicy(policyId)
	if err != nil {
		return "", err
	}

	return policyIdString, nil
}

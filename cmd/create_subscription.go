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
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/models"
	"intel/amber/tac/v1/validation"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// createSubscriptionCmd represents the createSubscription command
var createSubscriptionCmd = &cobra.Command{
	Use:   "subscription",
	Short: "Create a new subscription for a user",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("create subscription called")
		response, err := createSubscription(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Subscription: \n\n", response)
		return nil
	},
}

func init() {
	createCmd.AddCommand(createSubscriptionCmd)

	createSubscriptionCmd.Flags().StringP(constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	createSubscriptionCmd.Flags().StringP(constants.TenantIdParamName, "t", "", "Id of the tenant for whom the subscription needs to be created")
	createSubscriptionCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Amber service for which the subscription needs to be created")
	createSubscriptionCmd.Flags().StringP(constants.ProductIdParamName, "p", "", "Id of the Amber Product for which the subscription needs to be created")
	createSubscriptionCmd.Flags().StringP(constants.SubscriptionParamName, "d", "", "Description of the subscription that needs to be created")
	createSubscriptionCmd.MarkFlagRequired(constants.ApiKeyParamName)
	createSubscriptionCmd.MarkFlagRequired(constants.ServiceIdParamName)
	createSubscriptionCmd.MarkFlagRequired(constants.ProductIdParamName)
	createSubscriptionCmd.MarkFlagRequired(constants.SubscriptionParamName)
}

func createSubscription(cmd *cobra.Command) (string, error) {
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

	serviceIdString, err := cmd.Flags().GetString(constants.ServiceIdParamName)
	if err != nil {
		return "", err
	}

	serviceId, err := uuid.Parse(serviceIdString)
	if err != nil {
		return "", err
	}

	productIdString, err := cmd.Flags().GetString(constants.ProductIdParamName)
	if err != nil {
		return "", err
	}

	productId, err := uuid.Parse(productIdString)
	if err != nil {
		return "", err
	}

	subscriptionName, err := cmd.Flags().GetString(constants.SubscriptionParamName)
	if err != nil {
		return "", err
	}

	var subscriptionInfo = models.CreateSubscription{
		ProductId:   productId,
		Description: subscriptionName,
	}

	if err = validation.ValidateStrings([]string{subscriptionName}); err != nil {
		return "", errors.Wrap(err, "Invalid subscription name provided")
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, tenantId, apiKey)
	response, err := tmsClient.CreateSubscription(&subscriptionInfo, serviceId)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

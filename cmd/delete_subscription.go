/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 *
 *
 */

package cmd

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"net/http"
	"net/url"
	"time"
)

var deleteSubscriptionCmd = &cobra.Command{
	Use:   constants.SubscriptionCmd,
	Short: "Delete a subscription whose ID has been provided",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("delete subscription called")
		serviceId, err := deleteSubscription(cmd)
		if err != nil {
			return err
		}
		fmt.Printf("Deleted subscription with Id: %s \n\n", serviceId)
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteSubscriptionCmd)

	deleteSubscriptionCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	deleteSubscriptionCmd.Flags().StringP(constants.TenantIdParamName, "t", "", "Id of the tenant for whom the subscription needs to be created")
	deleteSubscriptionCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Amber service for which the subscription needs to be created")
	deleteSubscriptionCmd.Flags().StringP(constants.SubscriptionIdParamName, "d", "", "Id of the subscription which needs to be fetched (optional)")
	deleteSubscriptionCmd.MarkFlagRequired(constants.ApiKeyParamName)
	deleteSubscriptionCmd.MarkFlagRequired(constants.ServiceIdParamName)
	deleteSubscriptionCmd.MarkFlagRequired(constants.SubscriptionIdParamName)
}

func deleteSubscription(cmd *cobra.Command) (string, error) {
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

	subscriptionIdString, err := cmd.Flags().GetString(constants.SubscriptionIdParamName)
	if err != nil {
		return "", err
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, tenantId, apiKey)

	subscriptionId, err := uuid.Parse(subscriptionIdString)
	if err != nil {
		return "", err
	}

	err = tmsClient.DeleteSubscription(serviceId, subscriptionId)
	if err != nil {
		return "", err
	}

	return subscriptionIdString, nil
}

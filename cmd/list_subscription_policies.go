/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

var getSubscriptionPoliciesCmd = &cobra.Command{
	Use:   constants.PolicyCmd,
	Short: "Used to get the list of policy IDs linked to a subscription",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("list subscription policy called")
		response, err := getSubscriptionPolicies(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Services: \n\n", response)
		return nil
	},
}

func init() {
	getSubscriptionsCmd.AddCommand(getSubscriptionPoliciesCmd)

	getSubscriptionPoliciesCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	getSubscriptionPoliciesCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Amber service for which the subscription policies are to be fetched")
	getSubscriptionPoliciesCmd.Flags().StringP(constants.SubscriptionIdParamName, "s", "", "Id of the subscription for which the policies are to be fetched")
	getSubscriptionPoliciesCmd.MarkFlagRequired(constants.ApiKeyParamName)
	getSubscriptionPoliciesCmd.MarkFlagRequired(constants.ServiceIdParamName)
	getSubscriptionPoliciesCmd.MarkFlagRequired(constants.SubscriptionIdParamName)
}

func getSubscriptionPolicies(cmd *cobra.Command) (string, error) {
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
	subscriptionId, err := uuid.Parse(subscriptionIdString)
	if err != nil {
		return "", err
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, uuid.Nil, apiKey)
	response, err := tmsClient.GetSubscriptionPolicies(serviceId, subscriptionId)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

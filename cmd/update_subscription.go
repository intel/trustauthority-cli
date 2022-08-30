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
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/models"
	"intel/amber/tac/v1/validation"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// createSubscriptionCmd represents the createSubscription command
var updateSubscriptionCmd = &cobra.Command{
	Use:   constants.SubscriptionCmd,
	Short: "Update an existing subscription for a user",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("update subscription called")
		response, err := updateSubscription(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Subscription: \n\n", response)
		return nil
	},
}

func init() {
	updateCmd.AddCommand(updateSubscriptionCmd)

	updateSubscriptionCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	updateSubscriptionCmd.Flags().StringP(constants.TenantIdParamName, "t", "", "Id of the tenant for whom the subscription needs to be updated")
	updateSubscriptionCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Amber service for which the subscription needs to be updated")
	updateSubscriptionCmd.Flags().StringP(constants.ProductIdParamName, "p", "", "Id of the Amber Product for which the subscription needs to be updated")
	updateSubscriptionCmd.Flags().StringP(constants.SubscriptionDescriptionParamName, "d", "", "Description of the subscription that needs to be updated")
	updateSubscriptionCmd.Flags().StringP(constants.SubscriptionIdParamName, "u", "", "Id of the subscription that needs to be updated")
	updateSubscriptionCmd.Flags().StringSliceP(constants.PolicyIdsParamName, "i", []string{}, "List of comma separated policy IDs to be linked to the subscription")
	updateSubscriptionCmd.Flags().StringSliceP(constants.TagIdAndValuesParamName, "v", []string{}, "List of the comma separated tad Id and value pairs in the "+
		"following format:\n e03582e6-0709-42c2-a164-a687a970e040:Workload-AI,051800d0-cae5-48e7-8515-9801650fcd2b:60 V etc.")
	updateSubscriptionCmd.Flags().StringP(constants.ActivationStatus, "s", "", "Add activation status for subscription, should be \"active\" or \"inactive\"")
	updateSubscriptionCmd.MarkFlagRequired(constants.ApiKeyParamName)
	updateSubscriptionCmd.MarkFlagRequired(constants.ServiceIdParamName)
	updateSubscriptionCmd.MarkFlagRequired(constants.ProductIdParamName)
	updateSubscriptionCmd.MarkFlagRequired(constants.SubscriptionDescriptionParamName)
}

func updateSubscription(cmd *cobra.Command) (string, error) {
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

	subscriptionDescription, err := cmd.Flags().GetString(constants.SubscriptionDescriptionParamName)
	if err != nil {
		return "", err
	}

	subscriptionIdString, err := cmd.Flags().GetString(constants.SubscriptionIdParamName)
	if err != nil {
		return "", err
	}
	subscriptionId, err := uuid.Parse(subscriptionIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid subscription Id provided")
	}

	activationStatus, err := cmd.Flags().GetString(constants.ActivationStatus)
	if err != nil {
		return "", err
	} else if activationStatus != "active" && activationStatus != "inactive" {
		return "", errors.New("Activation status should be one of active or inactive")
	}

	policyIdsString, err := cmd.Flags().GetStringSlice(constants.PolicyIdsParamName)
	if err != nil {
		return "", err
	}

	var policyIds []uuid.UUID
	for _, policyId := range policyIdsString {
		policyUUID, err := uuid.Parse(policyId)
		if err != nil {
			return "", errors.Wrap(err, "Invalid policy ID found "+policyId+". Should be UUID.")
		}
		policyIds = append(policyIds, policyUUID)
	}

	tagIdValuesString, err := cmd.Flags().GetStringSlice(constants.TagIdAndValuesParamName)

	var tagIdValues []models.SubscriptionTagIdValue
	for _, tagIdValue := range tagIdValuesString {
		splitTag := strings.Split(tagIdValue, ":")
		if len(splitTag) != 2 {
			return "", errors.Wrap(err, "Tag Id value pairs are not provided in proper format, please check hel section for more details")
		}
		tagId, err := uuid.Parse(splitTag[0])
		if err != nil {
			return "", errors.Wrap(err, "Tag Id is not in proper format, should be UUID: "+tagId.String())
		}
		tagIdValues = append(tagIdValues, models.SubscriptionTagIdValue{TagId: tagId, Value: splitTag[1]})
	}

	var subscriptionInfo = models.UpdateSubscription{
		ProductId:    productId,
		Name:         subscriptionDescription,
		PolicyIds:    policyIds,
		TagIdsValues: tagIdValues,
		ServiceId:    serviceId,
		Status:       models.SubscriptionStatus(activationStatus),
	}

	if err = validation.ValidateStrings([]string{subscriptionDescription}); err != nil {
		return "", errors.Wrap(err, "Invalid subscription name provided")
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, tenantId, apiKey)
	response, err := tmsClient.UpdateSubscription(&subscriptionInfo, subscriptionId)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

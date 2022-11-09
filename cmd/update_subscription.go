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
	updateSubscriptionCmd.Flags().StringP(constants.SubscriptionNameParamName, "n", "", "Description of the subscription that needs to be updated")
	updateSubscriptionCmd.Flags().StringP(constants.SubscriptionIdParamName, "u", "", "Id of the subscription that needs to be updated")
	updateSubscriptionCmd.Flags().StringSliceP(constants.PolicyIdsParamName, "i", []string{}, "List of comma separated policy IDs to be linked to the subscription")
	updateSubscriptionCmd.Flags().StringSliceP(constants.TagKeyAndValuesParamName, "v", []string{}, "List of the comma separated tad Id and value pairs in the "+
		"following format:\n Workload:WorkloadAI,Workload:WorkloadEXE etc.")
	updateSubscriptionCmd.Flags().StringP(constants.SetExpiryDateParamName, "e", "", "Update the expiry date in the format yyyy-mm-dd for the new subscription")
	updateSubscriptionCmd.Flags().StringP(constants.ActivationStatus, "s", "", "Add activation status for subscription, should be one of \"Active\", \"Inactive\" or \"Cancelled\"")
	updateSubscriptionCmd.MarkFlagRequired(constants.ApiKeyParamName)
	updateSubscriptionCmd.MarkFlagRequired(constants.ServiceIdParamName)
	updateSubscriptionCmd.MarkFlagRequired(constants.ProductIdParamName)
	updateSubscriptionCmd.MarkFlagRequired(constants.SubscriptionNameParamName)
	updateSubscriptionCmd.MarkFlagRequired(constants.SubscriptionIdParamName)
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
		return "", errors.Wrap(err, "Invalid tenant id provided")
	}

	serviceIdString, err := cmd.Flags().GetString(constants.ServiceIdParamName)
	if err != nil {
		return "", err
	}

	serviceId, err := uuid.Parse(serviceIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid service id provided")
	}

	productIdString, err := cmd.Flags().GetString(constants.ProductIdParamName)
	if err != nil {
		return "", err
	}

	productId, err := uuid.Parse(productIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid product id provided")
	}

	subscriptionName, err := cmd.Flags().GetString(constants.SubscriptionNameParamName)
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
	} else if activationStatus != "" && activationStatus != constants.SubscriptionStatusActive &&
		activationStatus != constants.SubscriptionStatusInactive && activationStatus != constants.SubscriptionStatusCancelled {
		return "", errors.Errorf("Activation status should be one of %s, %s or %s", constants.SubscriptionStatusActive,
			constants.SubscriptionStatusInactive, constants.SubscriptionStatusCancelled)
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

	tagKeyValuesString, err := cmd.Flags().GetStringSlice(constants.TagKeyAndValuesParamName)

	var tagIdValues []models.SubscriptionTagIdValue
	for _, tagIdValue := range tagKeyValuesString {
		splitTag := strings.Split(tagIdValue, ":")
		if len(splitTag) != 2 {
			return "", errors.New("Tag Id value pairs are not provided in proper format, please check hel section for more details")
		}
		tagIdValues = append(tagIdValues, models.SubscriptionTagIdValue{Key: splitTag[0], Value: splitTag[1]})
	}

	expiryDateString, err := cmd.Flags().GetString(constants.SetExpiryDateParamName)
	if err != nil {
		return "", err
	}

	var subscriptionInfo = models.UpdateSubscription{
		ProductId:    productId,
		Name:         subscriptionName,
		PolicyIds:    policyIds,
		TagIdsValues: tagIdValues,
		ServiceId:    serviceId,
		Status:       models.SubscriptionStatus(activationStatus),
	}

	if expiryDateString != "" {
		date, err := time.Parse(constants.ExpiryDateInputFormat, expiryDateString)
		if err != nil {
			log.WithError(err).Error("Incorrect expiry date format provided")
			return "", errors.New("Incorrect expiry date provided. Date should be of the form yyyy-mm-dd")
		}
		subscriptionInfo.ExpiredAt = date
	}

	if err = validation.ValidateStrings([]string{subscriptionName}); err != nil {
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

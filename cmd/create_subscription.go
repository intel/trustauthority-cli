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
var createSubscriptionCmd = &cobra.Command{
	Use:   constants.SubscriptionCmd,
	Short: "Create a new subscription for a user",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("create subscription called")
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

	createSubscriptionCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	createSubscriptionCmd.Flags().StringP(constants.TenantIdParamName, "t", "", "Id of the tenant for whom the subscription needs to be created")
	createSubscriptionCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Amber service for which the subscription needs to be created")
	createSubscriptionCmd.Flags().StringP(constants.ProductIdParamName, "p", "", "Id of the Amber Product for which the subscription needs to be created")
	createSubscriptionCmd.Flags().StringP(constants.SubscriptionDescriptionParamName, "d", "", "Description of the subscription that needs to be created")
	createSubscriptionCmd.Flags().StringSliceP(constants.PolicyIdsParamName, "i", []string{}, "List of comma separated policy IDs to be linked to the subscription")
	createSubscriptionCmd.Flags().StringSliceP(constants.TagIdAndValuesParamName, "v", []string{}, "List of the comma separated tad Id and value pairs in the "+
		"following format:\n e03582e6-0709-42c2-a164-a687a970e040:Workload-AI,051800d0-cae5-48e7-8515-9801650fcd2b:60 V etc.")
	createSubscriptionCmd.Flags().StringP(constants.SetExpiryDateParamName, "e", "", "Set the expiry date in the format yyyy-mm-dd for the new subscription (default is 1 month)")
	createSubscriptionCmd.MarkFlagRequired(constants.ApiKeyParamName)
	createSubscriptionCmd.MarkFlagRequired(constants.ServiceIdParamName)
	createSubscriptionCmd.MarkFlagRequired(constants.ProductIdParamName)
	createSubscriptionCmd.MarkFlagRequired(constants.SubscriptionDescriptionParamName)
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

	subscriptionDescription, err := cmd.Flags().GetString(constants.SubscriptionDescriptionParamName)
	if err != nil {
		return "", err
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
	if err != nil {
		return "", err
	}

	var tagIdValues []models.SubscriptionTagIdValue
	for _, tagIdValue := range tagIdValuesString {
		splitTag := strings.Split(tagIdValue, ":")
		if len(splitTag) != 2 {
			return "", errors.New("Tag Id value pairs are not provided in proper format, please check help section for more details")
		}
		tagId, err := uuid.Parse(splitTag[0])
		if err != nil {
			return "", errors.Wrap(err, "Tag Id is not in proper format, should be UUID: "+splitTag[0])
		}
		tagIdValues = append(tagIdValues, models.SubscriptionTagIdValue{TagId: tagId, Value: splitTag[1]})
	}

	expiryDateString, err := cmd.Flags().GetString(constants.SetExpiryDateParamName)
	if err != nil {
		return "", err
	}

	var subscriptionInfo = models.CreateSubscription{
		ProductId:    productId,
		Description:  subscriptionDescription,
		PolicyIds:    policyIds,
		TagIdsValues: tagIdValues,
		CreatedBy:    tenantId,
		ServiceId:    serviceId,
		Status:       constants.SubscriptionStatusActive,
	}

	if expiryDateString == "" {
		fmt.Println("No expiry date provided. Setting subscription expiry date to 1 month from now")
		subscriptionInfo.ExpiredAt = time.Now().AddDate(0, 1, 0).UTC()
	} else {
		date, err := time.Parse(constants.ExpiryDateInputFormat, expiryDateString)
		if err != nil {
			log.WithError(err).Error("Incorrect expiry date format provided")
			return "", errors.New("Incorrect expiry date provided. Date should be of the form yyyy-mm-dd")
		}
		subscriptionInfo.ExpiredAt = date.UTC()
	}

	if err = validation.ValidateStrings([]string{subscriptionDescription}); err != nil {
		return "", errors.Wrap(err, "Invalid subscription name provided")
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, tenantId, apiKey)
	response, err := tmsClient.CreateSubscription(&subscriptionInfo)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

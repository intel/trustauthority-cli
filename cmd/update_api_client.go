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

var updateApiClientCmd = &cobra.Command{
	Use:   constants.ApiClientCmd,
	Short: "Update an existing api client for a user",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("update apiClient called")
		response, err := updateApiClient(cmd)
		if err != nil {
			return err
		}
		fmt.Println("ApiClient: \n\n", response)
		return nil
	},
}

func init() {
	updateCmd.AddCommand(updateApiClientCmd)

	updateApiClientCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	updateApiClientCmd.Flags().StringP(constants.TenantIdParamName, "t", "", "Id of the tenant for whom the api client needs to be updated")
	updateApiClientCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Amber service for which the api client needs to be updated")
	updateApiClientCmd.Flags().StringP(constants.ProductIdParamName, "p", "", "Id of the Amber Product for which the api client needs to be updated")
	updateApiClientCmd.Flags().StringP(constants.ApiClientNameParamName, "n", "", "Description of the api client that needs to be updated")
	updateApiClientCmd.Flags().StringP(constants.ApiClientIdParamName, "c", "", "Id of the api client that needs to be updated")
	updateApiClientCmd.Flags().StringSliceP(constants.PolicyIdsParamName, "i", []string{}, "List of comma separated policy IDs to be linked to the api client")
	updateApiClientCmd.Flags().StringSliceP(constants.TagKeyAndValuesParamName, "v", []string{}, "List of the comma separated tad Id and value pairs in the "+
		"following format:\n Workload:WorkloadAI,Workload:WorkloadEXE etc.")
	updateApiClientCmd.Flags().StringP(constants.SetExpiryDateParamName, "e", "", "Update the expiry date in the format yyyy-mm-dd for the new api client")
	updateApiClientCmd.Flags().StringP(constants.ActivationStatus, "s", "", "Add activation status for api client, should be one of \"Active\", \"Inactive\" or \"Cancelled\"")
	updateApiClientCmd.MarkFlagRequired(constants.ApiKeyParamName)
	updateApiClientCmd.MarkFlagRequired(constants.ServiceIdParamName)
	updateApiClientCmd.MarkFlagRequired(constants.ProductIdParamName)
	updateApiClientCmd.MarkFlagRequired(constants.ApiClientNameParamName)
	updateApiClientCmd.MarkFlagRequired(constants.ApiClientIdParamName)
}

func updateApiClient(cmd *cobra.Command) (string, error) {

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

	apiClientName, err := cmd.Flags().GetString(constants.ApiClientNameParamName)
	if err != nil {
		return "", err
	}

	apiClientIdString, err := cmd.Flags().GetString(constants.ApiClientIdParamName)
	if err != nil {
		return "", err
	}
	apiClientId, err := uuid.Parse(apiClientIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid api client Id provided")
	}

	activationStatus, err := cmd.Flags().GetString(constants.ActivationStatus)
	if err != nil {
		return "", err
	} else if activationStatus != "" && activationStatus != constants.ApiClientStatusActive &&
		activationStatus != constants.ApiClientStatusInactive && activationStatus != constants.ApiClientStatusCancelled {
		return "", errors.Errorf("Activation status should be one of %s, %s or %s", constants.ApiClientStatusActive,
			constants.ApiClientStatusInactive, constants.ApiClientStatusCancelled)
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

	var tagIdValues []models.ApiClientTagIdValue
	for _, tagIdValue := range tagKeyValuesString {
		splitTag := strings.Split(tagIdValue, ":")
		if len(splitTag) != 2 {
			return "", errors.New("Tag Id value pairs are not provided in proper format, please check hel section for more details")
		}
		tagIdValues = append(tagIdValues, models.ApiClientTagIdValue{Key: splitTag[0], Value: splitTag[1]})
	}

	expiryDateString, err := cmd.Flags().GetString(constants.SetExpiryDateParamName)
	if err != nil {
		return "", err
	}

	var apiClientInfo = models.UpdateApiClient{
		ProductId:    productId,
		Name:         apiClientName,
		PolicyIds:    policyIds,
		TagIdsValues: tagIdValues,
		ServiceId:    serviceId,
		Status:       models.ApiClientStatus(activationStatus),
	}

	if expiryDateString != "" {
		date, err := time.Parse(constants.ExpiryDateInputFormat, expiryDateString)
		if err != nil {
			log.WithError(err).Error("Incorrect expiry date format provided")
			return "", errors.New("Incorrect expiry date provided. Date should be of the form yyyy-mm-dd")
		}
		apiClientInfo.ExpiredAt = date
	}

	if err = validation.ValidateStrings([]string{apiClientName}); err != nil {
		return "", errors.Wrap(err, "Invalid api client name provided")
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, tenantId, apiKey)
	response, err := tmsClient.UpdateApiClient(&apiClientInfo, apiClientId)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

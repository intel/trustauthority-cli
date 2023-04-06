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

// createApiClientCmd represents the createApiClient command
var createApiClientCmd = &cobra.Command{
	Use:   constants.ApiClientCmd,
	Short: "Create a new api client for a user",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("create apiClient called")
		response, err := createApiClient(cmd)
		if err != nil {
			return err
		}
		fmt.Println("ApiClient: \n\n", response)
		return nil
	},
}

func init() {
	createCmd.AddCommand(createApiClientCmd)

	createApiClientCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Amber service for which the api client needs to be created")
	createApiClientCmd.Flags().StringP(constants.ProductIdParamName, "p", "", "Id of the Amber Product for which the api client needs to be created")
	createApiClientCmd.Flags().StringP(constants.ApiClientNameParamName, "n", "", "Name of the api client that needs to be created")
	createApiClientCmd.Flags().StringSliceP(constants.PolicyIdsParamName, "i", []string{}, "List of comma separated policy IDs to be linked to the api client")
	createApiClientCmd.Flags().StringSliceP(constants.TagKeyAndValuesParamName, "v", []string{}, "List of the comma separated tad Id and value pairs in the "+
		"following format:\n Workload:WorkloadAI,Workload:WorkloadEXE etc.")
	createApiClientCmd.MarkFlagRequired(constants.ServiceIdParamName)
	createApiClientCmd.MarkFlagRequired(constants.ProductIdParamName)
	createApiClientCmd.MarkFlagRequired(constants.ApiClientNameParamName)
}

func createApiClient(cmd *cobra.Command) (string, error) {

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
	if err != nil {
		return "", err
	}

	var tagKeyValues []models.ApiClientTagIdValue
	for _, tagIdValue := range tagKeyValuesString {
		splitTag := strings.Split(tagIdValue, ":")
		if len(splitTag) != 2 {
			return "", errors.New("Tag Id value pairs are not provided in proper format, please check help section for more details")
		}
		tagKeyValues = append(tagKeyValues, models.ApiClientTagIdValue{Key: splitTag[0], Value: splitTag[1]})
	}

	var apiClientInfo = models.CreateApiClient{
		ProductId:    productId,
		Name:         apiClientName,
		PolicyIds:    policyIds,
		TagIdsValues: tagKeyValues,
		ServiceId:    serviceId,
		Status:       constants.ApiClientStatusActive,
	}

	if err = validation.ValidateStrings([]string{apiClientName}); err != nil {
		return "", errors.Wrap(err, "Invalid api client name provided")
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, apiKey)
	response, err := tmsClient.CreateApiClient(&apiClientInfo)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

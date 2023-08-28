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
	"intel/tac/v1/client/tms"
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	models2 "intel/tac/v1/internal/models"
	"intel/tac/v1/models"
	"intel/tac/v1/utils"
	"intel/tac/v1/validation"
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
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Println("ApiClient: \n\n", response)
		fmt.Println("\nNOTE: There may be a delay of up to two (2) minutes before a new attestation API key is active.")
		fmt.Print("\n")
		return nil
	},
}

func init() {
	createCmd.AddCommand(createApiClientCmd)

	createApiClientCmd.Flags().StringP(constants.ServiceIdParamName, "r", "", "Id of the Trust Authority service for which the api client needs to be created")
	createApiClientCmd.Flags().StringP(constants.ProductIdParamName, "p", "", "Id of the Trust Authority Product for which the api client needs to be created")
	createApiClientCmd.Flags().StringP(constants.ApiClientNameParamName, "n", "", "Name of the api client that needs to be created")
	createApiClientCmd.Flags().StringSliceP(constants.PolicyIdsParamName, "i", []string{}, "List of comma separated policy IDs to be linked to the api client")
	createApiClientCmd.Flags().StringSliceP(constants.TagKeyAndValuesParamName, "v", []string{}, "List of the comma separated tad Id and value pairs in the "+
		"following format:\n Workload:WorkloadAI,Workload:WorkloadEXE etc.")
	createApiClientCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
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

	tmsUrl, err := url.Parse(configValues.TrustAuthorityBaseUrl + constants.TmsBaseUrl)
	if err != nil {
		return "", err
	}

	if err = setRequestId(cmd); err != nil {
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
	err = validation.ValidateApiClientName(apiClientName)
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
		if err = validation.ValidateTagName(splitTag[0]); err != nil {
			return "", err
		}
		if err = validation.ValidateTagValue(splitTag[1]); err != nil {
			return "", err
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

func setRequestId(cmd *cobra.Command) error {
	var err error
	models2.RespHeaderFields.RequestId, err = cmd.Flags().GetString(constants.RequestIdParamName)
	if err != nil {
		return err
	}

	if err = validation.ValidateRequestId(models2.RespHeaderFields.RequestId); err != nil {
		return err
	}
	return nil
}

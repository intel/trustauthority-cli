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
	"intel/tac/v1/utils"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// getProductsCmd represents the getProducts command
var getProductsCmd = &cobra.Command{
	Use:   constants.ProductCmd,
	Short: "Get list of products or a specific product under a specific service offer",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("list products called")
		response, err := getProducts(cmd)
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Println("Products: \n\n", response)
		return nil
	},
}

func init() {
	listCmd.AddCommand(getProductsCmd)

	getProductsCmd.Flags().StringP(constants.ServiceOfferIdParamName, "r", "", "Id of the Trust Authority  "+
		"service offer for which the product list needs to be fetched")
	getProductsCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
	getProductsCmd.MarkFlagRequired(constants.ServiceOfferIdParamName)
}

func getProducts(cmd *cobra.Command) (string, error) {
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

	serviceOfferIdString, err := cmd.Flags().GetString(constants.ServiceOfferIdParamName)
	if err != nil {
		return "", err
	}

	serviceOfferId, err := uuid.Parse(serviceOfferIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid service offer id provided")
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, apiKey)

	response, err := tmsClient.GetProducts(serviceOfferId)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

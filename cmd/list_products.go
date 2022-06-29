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

// getProductsCmd represents the getProducts command
var getProductsCmd = &cobra.Command{
	Use:   constants.ProductCmd,
	Short: "Get list of products or a specific product under a specific service offer",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("list products called")
		response, err := getProducts(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Products: \n\n", response)
		return nil
	},
}

func init() {
	listCmd.AddCommand(getProductsCmd)

	getProductsCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to "+
		"connect to amber services")
	getProductsCmd.Flags().StringP(constants.ServiceOfferIdParamName, "r", "", "Id of the Amber "+
		"service offer for which the product list needs to be fetched")
	getProductsCmd.MarkFlagRequired(constants.ApiKeyParamName)
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

	tmsUrl, err := url.Parse(configValues.AmberBaseUrl + constants.TmsBaseUrl)
	if err != nil {
		return "", err
	}

	serviceOfferIdString, err := cmd.Flags().GetString(constants.ServiceOfferIdParamName)
	if err != nil {
		return "", err
	}

	serviceOfferId, err := uuid.Parse(serviceOfferIdString)
	if err != nil {
		return "", err
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, uuid.Nil, apiKey)

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

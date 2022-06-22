/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// getServiceOffersCmd represents the getServiceOffers command
var getServiceOffersCmd = &cobra.Command{
	Use:   "serviceOffer",
	Short: "List all the service offers provided by Amber",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("get serviceOffers called")
		response, err := getServiceOffers()
		if err != nil {
			return err
		}
		fmt.Println("Service offers: \n\n", response)
		return nil
	},
}

func init() {
	listCmd.AddCommand(getServiceOffersCmd)

	getServiceOffersCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	getServiceOffersCmd.MarkFlagRequired(constants.ApiKeyParamName)
}

func getServiceOffers() (string, error) {
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

	tmsClient := tms.NewTmsClient(client, tmsUrl, uuid.Nil, apiKey)

	response, err := tmsClient.GetServiceOffers()
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

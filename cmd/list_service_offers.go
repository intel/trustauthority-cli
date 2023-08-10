/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/utils"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// getServiceOffersCmd represents the getServiceOffers command
var getServiceOffersCmd = &cobra.Command{
	Use:   constants.ServiceOfferCmd,
	Short: "List all the service offers provided by Amber",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("list serviceOffers called")
		response, err := getServiceOffers(cmd)
		if err != nil {
			return err
		}
		utils.PrintRequestAndTraceId()
		fmt.Println("Service offers: \n\n", response)
		return nil
	},
}

func init() {
	listCmd.AddCommand(getServiceOffersCmd)
	getServiceOffersCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
}

func getServiceOffers(cmd *cobra.Command) (string, error) {
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

	if err = setRequestId(cmd); err != nil {
		return "", err
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, apiKey)

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

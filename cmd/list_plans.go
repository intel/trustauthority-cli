/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 *
 *
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"net/http"
	"net/url"
	"time"
)

// getPlansCmd represents the getServices command
var getPlansCmd = &cobra.Command{
	Use:   constants.PlanCmd,
	Short: "Used to get the list of plans under a tenant",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("list plan called")
		response, err := getPlans(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Plans: \n\n", response)
		return nil
	},
}

func init() {
	listCmd.AddCommand(getPlansCmd)
	getPlansCmd.Flags().StringP(constants.ServiceOfferIdParamName, "r", "", "Id of the Amber service offer for which the plan needs to be fetched")
	getPlansCmd.Flags().StringP(constants.PlanIdParamName, "p", "", "Id of the Amber plan which needs to be fetched")
	getPlansCmd.MarkFlagRequired(constants.ServiceOfferIdParamName)
}

func getPlans(cmd *cobra.Command) (string, error) {
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
		return "", errors.Wrap(err, "Invalid service offer id provided")
	}

	planIdString, err := cmd.Flags().GetString(constants.PlanIdParamName)
	if err != nil {
		return "", err
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, uuid.Nil, apiKey)

	var responseBytes []byte
	if planIdString == "" {
		response, err := tmsClient.GetPlans(serviceOfferId)
		if err != nil {
			return "", err
		}

		responseBytes, err = json.MarshalIndent(response, "", "  ")
		if err != nil {
			return "", err
		}
	} else {
		planId, err := uuid.Parse(planIdString)
		if err != nil {
			return "", errors.Wrap(err, "Invalid plan id provided")
		}

		response, err := tmsClient.RetrievePlan(serviceOfferId, planId)
		if err != nil {
			return "", err
		}

		responseBytes, err = json.MarshalIndent(response, "", "  ")
		if err != nil {
			return "", err
		}
	}

	return string(responseBytes), nil
}

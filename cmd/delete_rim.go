/*
 * Copyright (C) 2026 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"fmt"
	"intel/tac/v1/client/pms"
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	"intel/tac/v1/utils"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// deleteRimCmd represents the deleteRim command
var deleteRimCmd = &cobra.Command{
	Use:   constants.RimCmd,
	Short: "Deletes a RIM",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("delete rim called")
		rimId, err := deleteRim(cmd)
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Printf("\nRIM %s deleted\n", rimId)
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteRimCmd)

	deleteRimCmd.Flags().StringP(constants.RimIdParamName, "r", "", "Id of the RIM to be deleted")
	deleteRimCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
	deleteRimCmd.MarkFlagRequired(constants.RimIdParamName)
}

func deleteRim(cmd *cobra.Command) (string, error) {
	configValues, err := config.LoadConfiguration()
	if err != nil {
		return "", err
	}
	client := &http.Client{
		Timeout: time.Duration(configValues.HTTPClientTimeout) * time.Second,
	}

	pmsUrl, err := url.Parse(configValues.TrustAuthorityBaseUrl + constants.PmsBaseUrl)
	if err != nil {
		return "", err
	}

	if err = setRequestId(cmd); err != nil {
		return "", err
	}

	rimIdString, err := cmd.Flags().GetString(constants.RimIdParamName)
	if err != nil {
		return "", err
	}

	rimId, err := uuid.Parse(rimIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid RIM id provided")
	}

	pmsClient := pms.NewPmsClient(client, pmsUrl, apiKey)

	err = pmsClient.DeleteRim(rimId)
	if err != nil {
		return "", err
	}

	return rimIdString, nil
}

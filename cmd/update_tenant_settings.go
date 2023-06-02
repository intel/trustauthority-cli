/*
 * Copyright (C) 2023 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/models"
	"intel/amber/tac/v1/validation"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// updateTenantSettingsCmd represents the update tenant-settings command
var updateTenantSettingsCmd = &cobra.Command{
	Use:   constants.TenantSettingsCmd,
	Short: "Update tenant settings for a tenant",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("update tenant settings command called")
		response, err := updateTenantSettings(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Updated Tenant Settings: \n", response)
		return nil
	},
}

func init() {
	updateCmd.AddCommand(updateTenantSettingsCmd)

	updateTenantSettingsCmd.Flags().StringP(constants.EmailIdParamName, "e", "", "The Email Id where the attestation failure notifications need to be sent")
	updateTenantSettingsCmd.Flags().BoolP(constants.DisableNotificationParamName, "d", false, "This parameter needs to be set to disable notification")
}

func updateTenantSettings(cmd *cobra.Command) (string, error) {
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

	emailId, err := cmd.Flags().GetString(constants.EmailIdParamName)
	if err != nil {
		return "", errors.Wrap(err, "Error fetching value of email-id parameter")
	}

	disableNotification, err := cmd.Flags().GetBool(constants.DisableNotificationParamName)
	if err != nil {
		return "", errors.Wrap(err, "Error fetching value of disable parameter")
	}

	tenantSettings := &models.AttestationFailureEmail{}
	if !disableNotification {
		err = validation.ValidateEmailAddress(emailId)
		if err != nil {
			return "", err
		}
		tenantSettings.AttestationFailureEmail = emailId
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, apiKey)
	response, err := tmsClient.UpdateTenantSettings(tenantSettings)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

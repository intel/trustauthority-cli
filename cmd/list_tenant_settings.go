/*
 * Copyright (C) 2023 Intel Corporation
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
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// listTenantSettingsCmd represents the list tenant-settings command
var listTenantSettingsCmd = &cobra.Command{
	Use:   constants.TenantSettingsCmd,
	Short: "List tenant settings for a tenant",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("list tenant settings command called")
		response, err := listTenantSettings(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Tenant Settings: \n", response)
		return nil
	},
}

func init() {
	listCmd.AddCommand(listTenantSettingsCmd)
}

func listTenantSettings(cmd *cobra.Command) (string, error) {
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

	tmsClient := tms.NewTmsClient(client, tmsUrl, apiKey)
	response, err := tmsClient.GetTenantSettings()
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

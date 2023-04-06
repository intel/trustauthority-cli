/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/models"
	"net/http"
	"net/url"
	"time"
)

var createTagCmd = &cobra.Command{
	Use:   constants.TagCmd,
	Short: "Creates a new tag",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("create tag called")
		response, err := createTag(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Tag: \n\n", response)
		return nil
	},
}

func init() {
	createCmd.AddCommand(createTagCmd)

	createTagCmd.Flags().StringP(constants.TagNameParamName, "n", "", "Name of the tag that needs to be created")
	createTagCmd.MarkFlagRequired(constants.TagNameParamName)
}

func createTag(cmd *cobra.Command) (string, error) {
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

	tagName, err := cmd.Flags().GetString(constants.TagNameParamName)
	if err != nil {
		return "", err
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, apiKey)

	createTagReq := &models.TagCreate{
		Name: tagName,
	}
	response, err := tmsClient.CreateTenantTag(createTagReq)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

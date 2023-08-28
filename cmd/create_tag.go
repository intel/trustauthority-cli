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
	"intel/tac/v1/client/tms"
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	"intel/tac/v1/models"
	"intel/tac/v1/utils"
	"intel/tac/v1/validation"
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
		utils.PrintRequestAndTraceId()
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
	createTagCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
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

	tmsUrl, err := url.Parse(configValues.TrustAuthorityBaseUrl + constants.TmsBaseUrl)
	if err != nil {
		return "", err
	}

	if err = setRequestId(cmd); err != nil {
		return "", err
	}

	tagName, err := cmd.Flags().GetString(constants.TagNameParamName)
	if err != nil {
		return "", err
	}
	if err = validation.ValidateTagName(tagName); err != nil {
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

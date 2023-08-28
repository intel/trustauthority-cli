/*
 * Copyright (C) 2023 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
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

// deleteTagCmd represents the deleteTag command
var deleteTagCmd = &cobra.Command{
	Use:   constants.TagCmd,
	Short: "Deletes a user-defined tag",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("delete tag called")
		tagId, err := deleteTag(cmd)
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Printf("\nTag %s deleted \n", tagId)
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteTagCmd)

	deleteTagCmd.Flags().StringP(constants.TagIdParamName, "t", "", "Id of the specific user defined tag which needs to be deleted")
	deleteTagCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
	deleteTagCmd.MarkFlagRequired(constants.TagIdParamName)
}

func deleteTag(cmd *cobra.Command) (string, error) {
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

	tagIdString, err := cmd.Flags().GetString(constants.TagIdParamName)
	if err != nil {
		return "", err
	}

	if tagIdString == "" {
		return "", errors.New("Tag Id cannot be empty")
	}

	tagId, err := uuid.Parse(tagIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid tag id provided")
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, apiKey)

	err = tmsClient.DeleteTenantTag(tagId)
	if err != nil {
		return "", err
	}

	return tagIdString, nil
}

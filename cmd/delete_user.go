/*
 * Copyright (C) 2022 Intel Corporation
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

// deleteUserCmd represents the deleteUser command
var deleteUserCmd = &cobra.Command{
	Use:   constants.UserCmd,
	Short: "Deletes a user under a tenant",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("delete user called")
		userId, err := deleteUser(cmd)
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Printf("\nUser %s deleted \n", userId)
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteUserCmd)
	deleteUserCmd.Flags().StringP(constants.UserIdParamName, "u", "", "Id of the specific user, the details for whom needs to be deleted")
	deleteUserCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
	deleteUserCmd.MarkFlagRequired(constants.UserIdParamName)
}

func deleteUser(cmd *cobra.Command) (string, error) {
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

	userIdString, err := cmd.Flags().GetString(constants.UserIdParamName)
	if err != nil {
		return "", err
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid user id provided")
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, apiKey)

	err = tmsClient.DeleteUser(userId)
	if err != nil {
		return "", err
	}

	return userIdString, nil
}

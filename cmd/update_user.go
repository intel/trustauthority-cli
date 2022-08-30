/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/models"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

var updateUserCmd = &cobra.Command{
	Use:   constants.UserCmd,
	Short: "Updates a user under a tenant",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("update user called")
		userId, err := updateUser(cmd)
		if err != nil {
			return err
		}
		fmt.Printf("\nUpdated User: %s \n\n", userId)
		return nil
	},
}

func init() {
	updateCmd.AddCommand(updateUserCmd)

	updateUserCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	updateUserCmd.Flags().StringP(constants.TenantIdParamName, "t", "", "Id of the tenant for whom the user needs to be created")
	updateUserCmd.Flags().StringP(constants.UserIdParamName, "u", "", "Id of the specific user")
	updateUserCmd.Flags().StringP(constants.EmailIdParamName, "e", "", "Email Id of the specific user to be updated")
	updateUserCmd.MarkFlagRequired(constants.ApiKeyParamName)
	updateUserCmd.MarkFlagRequired(constants.UserIdParamName)
	updateUserCmd.MarkFlagRequired(constants.UserRoleParamName)
}

func updateUser(cmd *cobra.Command) (string, error) {
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

	tenantIdString, err := cmd.Flags().GetString(constants.TenantIdParamName)
	if err != nil {
		return "", err
	}

	if tenantIdString == "" {
		tenantIdString = configValues.TenantId
	}

	tenantId, err := uuid.Parse(tenantIdString)
	if err != nil {
		return "", err
	}

	userIdString, err := cmd.Flags().GetString(constants.UserIdParamName)
	if err != nil {
		return "", err
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return "", err
	}

	emailId, err := cmd.Flags().GetString(constants.EmailIdParamName)
	if err != nil {
		return "", err
	}

	updateUserRoleReq := &models.UpdateUser{
		Id:    userId,
		Email: emailId,
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, tenantId, apiKey)

	response, err := tmsClient.UpdateUser(updateUserRoleReq)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

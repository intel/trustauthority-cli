/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"intel/tac/v1/client/tms"
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	"intel/tac/v1/models"
	"intel/tac/v1/utils"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

var (
	updateUserCmd = &cobra.Command{
		Use:   constants.UserCmd,
		Short: "Updates a user under a tenant",
		Long:  ``,
	}

	updateUserRoleCmd = &cobra.Command{
		Use:   constants.RoleCmd,
		Short: "Updates role of a user under a tenant",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("update user role called")
			userId, err := updateUserRole(cmd)
			if err != nil {
				return err
			}
			utils.PrintRequestAndTraceId()
			fmt.Printf("\nUpdated User: %s \n\n", userId)
			return nil
		},
	}
)

func init() {
	updateCmd.AddCommand(updateUserCmd)
	updateUserCmd.AddCommand(updateUserRoleCmd)

	updateUserRoleCmd.Flags().StringP(constants.UserIdParamName, "u", "", "Id of the specific user")
	updateUserRoleCmd.Flags().StringP(constants.UserRoleParamName, "r", "", "Role of the specific user that needs to be updated. Should be either Tenant Admin or User")
	updateUserRoleCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
	updateUserRoleCmd.MarkFlagRequired(constants.UserIdParamName)
	updateUserRoleCmd.MarkFlagRequired(constants.UserRoleParamName)
}

func updateUserRole(cmd *cobra.Command) (string, error) {
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

	userRole, err := cmd.Flags().GetString(constants.UserRoleParamName)
	if err != nil {
		return "", err
	}
	if userRole != constants.TenantAdminRole && userRole != constants.UserRole {
		return "", errors.Errorf("%s is not a valid user role. Roles should be either %s or %s", userRole,
			constants.TenantAdminRole, constants.UserRole)
	}

	updateUserRoleReq := &models.UpdateTenantUserRoles{
		UserId: userId,
		Role:   userRole,
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, apiKey)

	response, err := tmsClient.UpdateTenantUserRole(updateUserRoleReq)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

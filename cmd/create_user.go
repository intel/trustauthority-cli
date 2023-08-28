/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"intel/tac/v1/client/tms"
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	"intel/tac/v1/models"
	"intel/tac/v1/utils"
	"intel/tac/v1/validation"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// createUserCmd represents the createUser command
var createUserCmd = &cobra.Command{
	Use:   constants.UserCmd,
	Short: "Creates a new user under a tenant",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("create user called")
		response, err := createUser(cmd)
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Println("User: \n\n", response)
		return nil
	},
}

func init() {
	createCmd.AddCommand(createUserCmd)

	createUserCmd.Flags().StringP(constants.EmailIdParamName, "e", "", "Email id of the tenant user to be created")
	createUserCmd.Flags().StringP(constants.UserRoleParamName, "r", "", "Role of the tenant user to be created, should be one of Tenant Admin/User")
	createUserCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
	createUserCmd.MarkFlagRequired(constants.EmailIdParamName)
	createUserCmd.MarkFlagRequired(constants.UserRoleParamName)
}

func createUser(cmd *cobra.Command) (string, error) {
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

	emailId, err := cmd.Flags().GetString(constants.EmailIdParamName)
	if err != nil {
		return "", err
	}
	if err = validation.ValidateEmailAddress(emailId); err != nil {
		return "", err
	}

	userRole, err := cmd.Flags().GetString(constants.UserRoleParamName)
	if err != nil {
		return "", err
	}
	if userRole != constants.TenantAdminRole && userRole != constants.UserRole {
		return "", errors.Errorf("%s is not a valid user role. Roles should be either %s or %s", userRole,
			constants.TenantAdminRole, constants.UserRole)
	}

	var createUserInfo = &models.CreateTenantUser{
		Email: emailId,
		Role:  userRole,
	}

	tmsClient := tms.NewTmsClient(client, tmsUrl, apiKey)
	response, err := tmsClient.CreateUser(createUserInfo)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

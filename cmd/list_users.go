/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// getUsersCmd represents the getUsers command
var getUsersCmd = &cobra.Command{
	Use:   "user",
	Short: "Get user(s) under tenant",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("get users called")
		response, err := getUsers(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Users: \n\n", response)
		return nil
	},
}

func init() {
	listCmd.AddCommand(getUsersCmd)

	getUsersCmd.Flags().StringVarP(&apiKey, constants.ApiKeyParamName, "a", "", "API key to be used to connect to amber services")
	getUsersCmd.Flags().StringP(constants.TenantIdParamName, "t", "", "Id of the tenant for whom the subscription needs to be created")
	getUsersCmd.Flags().StringP(constants.UserIdParamName, "u", "", "Id of the specific user the details for whom needs to be fetched")
	getUsersCmd.MarkFlagRequired(constants.ApiKeyParamName)

}

func getUsers(cmd *cobra.Command) (string, error) {
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

	tmsClient := tms.NewTmsClient(client, tmsUrl, tenantId, apiKey)

	var responseBytes []byte
	if userIdString == "" {
		response, err := tmsClient.GetUsers()
		if err != nil {
			return "", err
		}

		responseBytes, err = json.MarshalIndent(response, "", "  ")
		if err != nil {
			return "", err
		}
	} else {
		userId, err := uuid.Parse(userIdString)
		if err != nil {
			return "", err
		}

		response, err := tmsClient.RetrieveUser(userId)
		if err != nil {
			return "", err
		}

		responseBytes, err = json.MarshalIndent(response, "", "  ")
		if err != nil {
			return "", err
		}
	}

	return string(responseBytes), nil
}

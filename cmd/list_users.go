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
	"intel/amber/tac/v1/client/tms"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/validation"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

// getUsersCmd represents the getUsers command
var getUsersCmd = &cobra.Command{
	Use:   constants.UserCmd,
	Short: "Get user(s) under tenant",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("list users called")
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

	getUsersCmd.Flags().StringP(constants.EmailIdParamName, "e", "", "Email Id of the Tenant User to be retrieved")
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

	tmsClient := tms.NewTmsClient(client, tmsUrl, uuid.Nil, apiKey)

	var responseBytes []byte
	response, err := tmsClient.GetUsers()
	if err != nil {
		return "", err
	}

	emailIdString, err := cmd.Flags().GetString(constants.EmailIdParamName)
	if err != nil {
		return "", err
	}

	if emailIdString != "" {
		err := validation.ValidateEmailAddress(emailIdString)
		if err != nil {
			return "", err
		}
		for _, user := range response {
			if user.Email == emailIdString {
				responseBytes, err = json.MarshalIndent(user, "", "  ")
				if err != nil {
					return "", err
				}
				return string(responseBytes), nil
			}
		}
		return "", errors.New("User associated with the email Id provided in input was not found")
	} else {
		fmt.Println("Email ID was not provided, listing all users....")
		responseBytes, err = json.MarshalIndent(response, "", "  ")
		if err != nil {
			return "", err
		}

		return string(responseBytes), nil
	}
}

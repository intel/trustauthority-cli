/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	"intel/tac/v1/internal/models"
	"intel/tac/v1/utils"
	"intel/tac/v1/validation"
	"os"
	"path/filepath"
)

var (
	apiKey string
)

// tenantCmd represents the base command when called without any subcommands
var tenantCmd = &cobra.Command{
	Use:   constants.RootCmd,
	Short: "Intel Trust Authority CLI used to run the tasks for tenant admin/user",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the tenantCmd.
func Execute() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error fetching user home directory path. Error: ", err.Error())
	}

	cleanedLogPath := filepath.Clean(userHomeDir + constants.LogFilePath)

	logFile, err := os.OpenFile(cleanedLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, constants.DefaultFilePermission)
	if err != nil {
		fmt.Println("Error opening/creating log file: " + err.Error())
		os.Exit(1)
	}
	tenantCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		configValues, err := config.LoadConfiguration()
		if err != nil {
			if logErr := utils.SetUpLogs(logFile, constants.DefaultLogLevel); logErr != nil {
				return logErr
			}
			logrus.WithError(err).Error("Error loading configuration")
			return err
		} else {
			if err := utils.SetUpLogs(logFile, configValues.LogLevel); err != nil {
				return err
			}
		}

		//API key is not needed for generating policy JWT or setting up config, API key check is skipped for these commands
		cmdListWithNoApiKey := map[string]bool{constants.PolicyJwtCmd: true, constants.SetupConfigCmd: true,
			constants.UninstallCmd: true, constants.VersionCmd: true}
		//API key is not needed for generating policy JWT or setting up config, API key check is skipped for these 2 commands
		if ok := cmdListWithNoApiKey[cmd.Name()]; !ok {
			apiKey = configValues.TrustAuthorityApiKey
			if err = validation.ValidateTrustAuthorityAPIKey(configValues.TrustAuthorityApiKey); err != nil {
				// check if jwt token is passed instead of api-key (packaged software use-case)
				if err = validation.ValidateTrustAuthorityJwt(configValues.TrustAuthorityApiKey); err != nil {
					return errors.New("Invalid Trust Authority Api key, API key should be a base64 encoded string or a JWT")
				}
			}
			if configValues.TrustAuthorityBaseUrl != "" {
				err = validation.ValidateURL(configValues.TrustAuthorityBaseUrl)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	err = tenantCmd.Execute()
	if err != nil {
		//Need to set it here separately as well since previously we are setting it only for the executed command
		logrus.SetOutput(logFile)
		logrus.WithField(constants.HTTPHeaderKeyRequestId, models.RespHeaderFields.RequestId).
			WithField(constants.HTTPHeaderKeyTraceId, models.RespHeaderFields.TraceId).Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}

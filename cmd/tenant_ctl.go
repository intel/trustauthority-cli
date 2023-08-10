/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/internal/models"
	"intel/amber/tac/v1/utils"
	"intel/amber/tac/v1/validation"
	"os"
	"path/filepath"
)

var (
	apiKey string
)

// tenantCmd represents the base command when called without any subcommands
var tenantCmd = &cobra.Command{
	Use:   constants.RootCmd,
	Short: "Tenant CLI used to run the tasks for tenant admin/user",
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
			constants.UninstallCmd: true}
		//API key is not needed for generating policy JWT or setting up config, API key check is skipped for these 2 commands
		if ok := cmdListWithNoApiKey[cmd.Name()]; !ok {
			apiKey = configValues.AmberApiKey
			if err := validation.ValidateAmberAPIKey(apiKey); err != nil {
				return err
			}
		}
		return nil
	}
	err = tenantCmd.Execute()
	if err != nil {
		logrus.WithField(constants.HTTPHeaderKeyRequestId, models.RespHeaderFields.RequestId).
			WithField(constants.HTTPHeaderKeyTraceId, models.RespHeaderFields.TraceId).Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}

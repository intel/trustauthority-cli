/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// setupConfigCmd represents the setup command
var (
	setupConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "Setup configuration for Amber CLI",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("setup called")
			if err := setupConfig(); err != nil {
				return err
			}
			return nil
		},
	}

	envFilePath string
)

func init() {
	tenantCmd.AddCommand(setupConfigCmd)

	setupConfigCmd.Flags().StringVarP(&envFilePath, constants.EnvFileParamName, "v", "", "Path for the env file to be used to update the configuration")
	setupConfigCmd.MarkFlagRequired(constants.EnvFileParamName)
}

func setupConfig() error {
	var err error

	if _, err = os.Stat(constants.DefaultConfigFilePath); err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(constants.DefaultConfigFilePath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if err = utils.ReadAnswerFileToEnv(envFilePath); err != nil {
		return err
	}

	//set default
	viper.SetDefault(constants.Loglevel, constants.DefaultLogLevel)
	viper.SetDefault(constants.HttpClientTimeout, constants.DefaultHttpClientTimeout)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	configValues := &config.Configuration{}

	configValues.AmberBaseUrl = viper.GetString(constants.AmberBaseUrl)
	tenantId, err := uuid.Parse(viper.GetString(constants.TenantId))
	if err != nil {
		return err
	}
	configValues.TenantId = tenantId.String()
	configValues.LogLevel, err = log.ParseLevel(viper.GetString(constants.Loglevel))
	if err != nil {
		return err
	}
	configValues.HTTPClientTimeout = viper.GetInt(constants.HttpClientTimeout)

	if configValues.AmberBaseUrl == "" {
		fmt.Println("Amber base URL needs to be provided in configuration")
		os.Exit(1)
	}

	if err = configValues.Save(constants.DefaultConfigFilePath); err != nil {
		return err
	}
	return nil
}

/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
)

// setupConfigCmd represents the setup command
var (
	setupConfigCmd = &cobra.Command{
		Use:   constants.SetupConfigCmd,
		Short: "Setup configuration for Amber CLI",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("setup called")
			if err := config.SetupConfig(envFilePath); err != nil {
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

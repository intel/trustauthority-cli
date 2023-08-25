/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"intel/tac/v1/constants"
	"intel/tac/v1/utils"
)

var versionCmd = &cobra.Command{
	Use:   constants.VersionCmd,
	Short: "Get version of Intel Trust Authority CLI",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := utils.GetVersion()
		if err != nil {
			return err
		}
		fmt.Printf("\nVersion: %s\n\n", version)
		return nil
	},
}

func init() {
	tenantCmd.AddCommand(versionCmd)
}

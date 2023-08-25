/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"github.com/spf13/cobra"
	"intel/tac/v1/constants"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   constants.UpdateCmd,
	Short: "Updates a resource",
	Long:  ``,
}

func init() {
	tenantCmd.AddCommand(updateCmd)
}

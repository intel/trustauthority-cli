/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"github.com/spf13/cobra"
	"intel/amber/tac/v1/constants"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   constants.ListCmd,
	Short: "Lists a resource or group of resources",
	Long:  ``,
}

func init() {
	tenantCmd.AddCommand(listCmd)
}

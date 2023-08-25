/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"github.com/spf13/cobra"
	"intel/tac/v1/constants"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   constants.DeleteCmd,
	Short: "Delete a resource",
	Long:  ``,
}

func init() {
	tenantCmd.AddCommand(deleteCmd)
}

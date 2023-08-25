/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"github.com/spf13/cobra"
	"intel/tac/v1/constants"
)

var createCmd = &cobra.Command{
	Use:   constants.CreateCmd,
	Short: "Create a resource",
	Long:  ``,
}

func init() {
	tenantCmd.AddCommand(createCmd)
}

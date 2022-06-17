/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	apiKey string
)

// tenantCmd represents the base command when called without any subcommands
var tenantCmd = &cobra.Command{
	Use:   "tenantctl",
	Short: "Tenant CLI used to run the tasks for tenant admin/user",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the tenantCmd.
func Execute() {
	err := tenantCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}

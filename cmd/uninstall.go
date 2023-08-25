/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"intel/tac/v1/constants"
	"os"

	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   constants.UninstallCmd,
	Short: "Uninstall Intel Trust Authority CLI",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Uninstalling Intel Trust Authority CLI")
		if err := uninstall(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	tenantCmd.AddCommand(uninstallCmd)
}

func uninstall() error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error fetching user home directory path. Error: ", err.Error())
	}

	log.Info("removing : ", userHomeDir+constants.BinDir)
	err = os.RemoveAll(userHomeDir + constants.BinDir)
	if err != nil {
		log.WithError(err).Error("Error removing home dir: ", userHomeDir+constants.BinDir)
	}

	log.Info("removing : ", userHomeDir+constants.ConfigDir)
	err = os.RemoveAll(userHomeDir + constants.ConfigDir)
	if err != nil {
		log.WithError(err).Error("Error removing config dir: ", userHomeDir+constants.ConfigDir)
	}

	log.Info("removing : ", userHomeDir+constants.LogDir)
	err = os.RemoveAll(userHomeDir + constants.LogDir)
	if err != nil {
		log.WithError(err).Error("Error removing log dir: ", userHomeDir+constants.LogDir)
	}

	fmt.Println("Intel Trust Authority CLI uninstalled successfully")
	return nil
}

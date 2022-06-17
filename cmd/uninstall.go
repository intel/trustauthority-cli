/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"intel/amber/tac/v1/constants"
	"os"

	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Tenant CLI",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Uninstalling Tenant CLI")
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
	fmt.Println("removing : ", constants.HomeDir)
	err := os.RemoveAll(constants.HomeDir)
	if err != nil {
		log.WithError(err).Error("Error removing home dir: ", constants.HomeDir)
	}

	fmt.Println("removing : ", constants.ConfigDir)
	err = os.RemoveAll(constants.ConfigDir)
	if err != nil {
		log.WithError(err).Error("Error removing config dir: ", constants.ConfigDir)
	}

	fmt.Println("Tenant CLI uninstalled successfully")
	return nil
}

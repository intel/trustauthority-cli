/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/test"
	"os"
	"strings"
	"testing"
)

var tempConfigFile *os.File
var err error

func init() {
	tempDir := os.TempDir()
	tempConfigFile, err = os.Create(tempDir + "/" + constants.ConfigFileName + "." + constants.ConfigFileExtension)
	if err != nil {
		fmt.Println(err.Error())
	}
	_ = os.Rename(tempConfigFile.Name(), strings.TrimRight(tempConfigFile.Name(), "1234567890"))
	viper.SetConfigName(constants.ConfigFileName)
	viper.SetConfigType(constants.ConfigFileExtension)
	viper.AddConfigPath(tempDir)
}

func TestCreatePolicyCmd(t *testing.T) {
	server := test.MockServer(t)
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args: []string{constants.CreateCmd, constants.PolicyCmd, "-n", "Sample_Policy_SGX", "-t", "Appraisal policy",
				"-r", "e8a72b7e-c4b1-4bdc-bf40-68f23c68a2aa", "-a", "SGX Attestation", "-f", "../test/resources/rego-policy.txt"},
			wantErr:     false,
			description: "Test Create Policy",
		},
		{
			args: []string{constants.CreateCmd, constants.PolicyCmd, "-n", "Sample_Policy_SGX", "-t", "Appraisal policy",
				"-r", "e8a72b7e-c4b1-4bdc-bf40-68f23c68a2aa", "-a", "SGX Attestation", "-f", ""},
			wantErr:     true,
			description: "Test Create Policy With Invalid File Path",
		},
	}

	createCmd.AddCommand(createPolicyCmd)
	tenantCmd.AddCommand(createCmd)

	for _, tc := range tt {
		_, err := execute(t, tenantCmd, tc.args)

		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func execute(t *testing.T, c *cobra.Command, args []string) (string, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	c.SetOut(buf)
	c.SetErr(buf)
	c.SetArgs(args)

	err := c.Execute()
	return strings.TrimSpace(buf.String()), err
}

/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/test"
	"testing"
)

func TestDeletePolicyCommandWithInvalidUrl(t *testing.T) {
	test.SetupMockConfiguration("invalid url", tempConfigFile)
	load, err := config.LoadConfiguration()
	assert.NoError(t, err)
	viper.Set("amber-base-url", "bogus\nbase\nURL")

	invalidUrlTc := struct {
		args        []string
		wantErr     bool
		description string
	}{
		args:        []string{constants.DeleteCmd, constants.PolicyCmd, "-p", "e48dabc5-9608-4ff3-aaed-f25909ab9de1"},
		wantErr:     true,
		description: "Test delete policy using invalid URL",
	}
	deleteCmd.AddCommand(deletePolicyCmd)
	tenantCmd.AddCommand(deleteCmd)

	_, err = execute(t, tenantCmd, invalidUrlTc.args)
	viper.Set("amber-base-url", load.AmberBaseUrl)
	assert.Error(t, err)
}

func TestDeletePolicyCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.DeleteCmd, constants.PolicyCmd, "-p", "e48dabc5-9608-4ff3-aaed-f25909ab9de1"},
			wantErr: false,
		},
		{
			args:        []string{constants.DeleteCmd, constants.PolicyCmd, "-p", "invalid id"},
			wantErr:     true,
			description: "Test Invalid policy id provided",
		},
	}

	deleteCmd.AddCommand(deletePolicyCmd)
	tenantCmd.AddCommand(deleteCmd)

	for _, tc := range tt {
		_, err := execute(t, tenantCmd, tc.args)

		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

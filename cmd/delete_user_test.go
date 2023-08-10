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

func TestDeleteUserCommandWithInvalidUrl(t *testing.T) {
	test.SetupMockConfiguration("invalid url", tempConfigFile)
	load, err := config.LoadConfiguration()
	assert.NoError(t, err)
	viper.Set("amber-base-url", "bogus\nbase\nURL")

	invalidUrlTc := struct {
		args        []string
		wantErr     bool
		description string
	}{
		args:        []string{constants.DeleteCmd, constants.UserCmd, "-u", "23011406-6f3b-4431-9363-4e1af9af6b13"},
		wantErr:     true,
		description: "Test delete user using invalid URL",
	}
	deleteCmd.AddCommand(deleteUserCmd)
	tenantCmd.AddCommand(deleteCmd)

	_, err = execute(t, tenantCmd, invalidUrlTc.args)
	viper.Set("amber-base-url", load.AmberBaseUrl)
	assert.Error(t, err)
}

func TestDeleteUserCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.DeleteCmd, constants.UserCmd, "-q", "valid-id", "-u", "23011406-6f3b-4431-9363-4e1af9af6b13"},
			wantErr: false,
		},
		{
			args:        []string{constants.DeleteCmd, constants.UserCmd, "-u", "invalid id"},
			wantErr:     true,
			description: "Test Invalid user id provided",
		},
		{
			args:        []string{constants.DeleteCmd, constants.UserCmd, "-q", "@#$invalid-id", "-u", "23011406-6f3b-4431-9363-4e1af9af6b13"},
			wantErr:     true,
			description: "Test Invalid request id provided",
		},
	}

	deleteCmd.AddCommand(deleteUserCmd)
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

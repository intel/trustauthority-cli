/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	"intel/tac/v1/test"
	"testing"
)

func TestUpdateUserRoleCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args: []string{constants.UpdateCmd, constants.UserCmd, constants.RoleCmd,
				"-r", "", "-u", "23011406-6f3b-4431-9363-4e1af9af6b13"},
			wantErr:     true,
			description: "Test Invalid user role. Roles should be either Tenant Admin or User",
		},
		{
			args: []string{constants.UpdateCmd, constants.UserCmd, constants.RoleCmd, "-q", "valid-id",
				"-r", "User", "-u", "23011406-6f3b-4431-9363-4e1af9af6b13"},
			wantErr: false,
		},
		{
			args: []string{constants.UpdateCmd, constants.UserCmd, constants.RoleCmd,
				"-r", "User", "-u", "invalid id"},
			wantErr:     true,
			description: "Test Invalid user id provided",
		},
		{
			args: []string{constants.UpdateCmd, constants.UserCmd, constants.RoleCmd, "-q", "@#$invalid-id",
				"-r", "User", "-u", "23011406-6f3b-4431-9363-4e1af9af6b13"},
			wantErr:     true,
			description: "Test Invalid request id provided",
		},
	}

	updateUserCmd.AddCommand(updateUserRoleCmd)
	updateCmd.AddCommand(updateUserCmd)
	tenantCmd.AddCommand(updateCmd)

	for _, tc := range tt {
		_, err := execute(t, tenantCmd, tc.args)

		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestUpdateUserRoleCommandInvalidUrl(t *testing.T) {
	test.SetupMockConfiguration("invalid url", tempConfigFile)
	load, err := config.LoadConfiguration()
	assert.NoError(t, err)

	invalidUrlTc := []struct {
		args        []string
		wantErr     bool
		url         string
		description string
	}{
		{
			args: []string{constants.UpdateCmd, constants.UserCmd, constants.RoleCmd, "-q", "valid-id",
				"-r", "User", "-u", "23011406-6f3b-4431-9363-4e1af9af6b13"},
			wantErr:     true,
			url:         "bogus\nbase\nURL",
			description: "Test update user role using invalid URL",
		},

		{
			args: []string{constants.UpdateCmd, constants.UserCmd, constants.RoleCmd, "-q", "valid-id",
				"-r", "User", "-u", "23011406-6f3b-4431-9363-4e1af9af6b13"},
			wantErr:     true,
			url:         "a/b/c",
			description: "Invalid send request provided for update user role command ",
		},
	}
	updateUserCmd.AddCommand(updateUserRoleCmd)
	updateCmd.AddCommand(updateUserCmd)
	tenantCmd.AddCommand(updateCmd)

	for _, tc := range invalidUrlTc {
		viper.Set("trustauthority-url", tc.url)
		_, err := execute(t, tenantCmd, tc.args)
		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
	viper.Set("trustauthority-url", load.TrustAuthorityBaseUrl)
}

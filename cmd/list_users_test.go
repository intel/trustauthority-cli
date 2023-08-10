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

func TestListUsersCommandWithInvalidUrl(t *testing.T) {
	test.SetupMockConfiguration("invalid url", tempConfigFile)
	load, err := config.LoadConfiguration()
	assert.NoError(t, err)
	viper.Set("amber-base-url", "bogus\nbase\nURL")

	invalidUrlTc := struct {
		args        []string
		wantErr     bool
		description string
	}{
		args:        []string{constants.ListCmd, constants.UserCmd},
		wantErr:     false,
		description: "Test list users using invalid URL",
	}
	listCmd.AddCommand(getUsersCmd)
	tenantCmd.AddCommand(listCmd)

	_, err = execute(t, tenantCmd, invalidUrlTc.args)
	viper.Set("amber-base-url", load.AmberBaseUrl)
	assert.Error(t, err)
}
func TestListUsersCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.ListCmd, constants.UserCmd, "-q", "valid-id"},
			wantErr: false,
		},
		{
			args:    []string{constants.ListCmd, constants.UserCmd, "-e", "notfound@mail.com"},
			wantErr: true,
		},
		{
			args:    []string{constants.ListCmd, constants.UserCmd, "-e", "arijitgh@gmail.com"},
			wantErr: false,
		},
		{
			args:    []string{constants.ListCmd, constants.UserCmd, "-q", "@#$invalid-id"},
			wantErr: true,
		},
	}

	listCmd.AddCommand(getUsersCmd)
	tenantCmd.AddCommand(listCmd)

	for _, tc := range tt {
		_, err := execute(t, tenantCmd, tc.args)

		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

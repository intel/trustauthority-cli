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

func TestCreateUserCommandWithInvalidUrl(t *testing.T) {
	test.SetupMockConfiguration("invalid url", tempConfigFile)
	load, err := config.LoadConfiguration()
	assert.NoError(t, err)
	viper.Set("amber-base-url", "bogus\nbase\nURL")

	invalidUrlTc := struct {
		args        []string
		wantErr     bool
		description string
	}{
		args:        []string{constants.CreateCmd, constants.UserCmd, "-e", "test@mail.com", "-r", "User"},
		wantErr:     true,
		description: "Test Create user using invalid URL",
	}
	createCmd.AddCommand(createUserCmd)
	tenantCmd.AddCommand(createCmd)

	_, err = execute(t, tenantCmd, invalidUrlTc.args)
	viper.Set("amber-base-url", load.AmberBaseUrl)
	assert.Error(t, err)
}

func TestCreateUserCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:        []string{constants.CreateCmd, constants.UserCmd, "-e", "test@mail.com", "-r", "User"},
			wantErr:     false,
			description: "Create user",
		},
		{
			args:        []string{constants.CreateCmd, constants.UserCmd, "-e", "#!bash#script", "-r", "User"},
			wantErr:     true,
			description: "Create user using invalid email id",
		},
		{
			args:        []string{constants.CreateCmd, constants.UserCmd, "-e", "test@mail.com", "-r", "Administrator"},
			wantErr:     true,
			description: "Create user using invalid role",
		},
	}

	createCmd.AddCommand(createUserCmd)
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

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

func TestCreateUserCommandWithInvalidUrl(t *testing.T) {
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
			args:        []string{constants.CreateCmd, constants.UserCmd, "-e", "test@mail.com", "-r", "User"},
			wantErr:     true,
			url:         "bogus\nbase\nURL",
			description: "Test Create user using invalid URL",
		},
		{
			args:        []string{constants.CreateCmd, constants.UserCmd, "-e", "test@mail.com", "-r", "User"},
			wantErr:     true,
			url:         "a/b/c",
			description: "Invalid send request provided for create user command",
		},
	}
	createCmd.AddCommand(createUserCmd)
	tenantCmd.AddCommand(createCmd)

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
			args:        []string{constants.CreateCmd, constants.UserCmd, "-q", "valid-id", "-e", "test@mail.com", "-r", "User"},
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
		{
			args:        []string{constants.CreateCmd, constants.UserCmd, "-q", "@#$invalid-id", "-e", "test@mail.com", "-r", "User"},
			wantErr:     true,
			description: "Create user using invalid request ID",
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

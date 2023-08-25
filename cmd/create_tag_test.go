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

func TestCreateTagWithInvalidUrl(t *testing.T) {
	test.SetupMockConfiguration("invalid url", tempConfigFile)
	load, err := config.LoadConfiguration()
	assert.NoError(t, err)
	viper.Set("trustauthority-url", "bogus\nbase\nURL")

	invalidUrlTc := struct {
		args        []string
		wantErr     bool
		description string
	}{
		args:        []string{constants.CreateCmd, constants.TagCmd, "-n", "Test Tag"},
		wantErr:     true,
		description: "Test Create tag using invalid URL",
	}

	createCmd.AddCommand(createTagCmd)
	tenantCmd.AddCommand(createCmd)

	_, err = execute(t, tenantCmd, invalidUrlTc.args)
	viper.Set("trustauthority-url", load.TrustAuthorityBaseUrl)
	assert.Error(t, err)
}

func TestCreateTagCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:        []string{constants.CreateCmd, constants.TagCmd, "-q", "valid-id", "-n", "Test_Tag"},
			wantErr:     false,
			description: "Create a tag",
		},
		{
			args:        []string{constants.CreateCmd, constants.TagCmd, "-n", "Test @Tag"},
			wantErr:     true,
			description: "Create a tag using invalid tag name",
		},
		{
			args:        []string{constants.CreateCmd, constants.TagCmd, "-q", "@#$invalid-id", "-n", "Test_Tag"},
			wantErr:     true,
			description: "Create a tag using invalid request ID",
		},
	}

	createCmd.AddCommand(createTagCmd)
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

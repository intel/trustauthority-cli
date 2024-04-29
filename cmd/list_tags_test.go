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

func TestListTagCommandWithInvalidUrl(t *testing.T) {
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
			args:        []string{constants.ListCmd, constants.TagCmd},
			wantErr:     true,
			url:         "bogus\nbase\nURL",
			description: "Test list tags using invalid URL",
		},
		{
			args:        []string{constants.ListCmd, constants.TagCmd, "-q", "valid-id"},
			wantErr:     true,
			url:         "a/b/c",
			description: "Invalid send request provided for list tags command ",
		},
	}
	listCmd.AddCommand(listTagCmd)
	tenantCmd.AddCommand(listCmd)

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

func TestListTagCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.ListCmd, constants.TagCmd, "-q", "valid-id"},
			wantErr: false,
		},
		{
			args:    []string{constants.ListCmd, constants.TagCmd, "-q", "@#$invalid-id"},
			wantErr: true,
		},
	}

	listCmd.AddCommand(listTagCmd)
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

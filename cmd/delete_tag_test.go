/*
 * Copyright (C) 2023 Intel Corporation
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

func TestDeleteTagCommandWithInvalidUrl(t *testing.T) {
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
			args:        []string{constants.DeleteCmd, constants.TagCmd, "-t", "23011406-6f3b-4431-9363-4e1af9af6b13"},
			wantErr:     true,
			url:         "bogus\nbase\nURL",
			description: "Test delete tag using invalid URL",
		},
		{
			args:        []string{constants.DeleteCmd, constants.TagCmd, "-t", "23011406-6f3b-4431-9363-4e1af9af6b13"},
			wantErr:     true,
			url:         "a/b/c",
			description: "Invalid send request provided for delete tag command",
		},
	}

	deleteCmd.AddCommand(deleteTagCmd)
	tenantCmd.AddCommand(deleteCmd)

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
func TestDeleteTagCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.DeleteCmd, constants.TagCmd, "-q", "valid-id", "-t", "23011406-6f3b-4431-9363-4e1af9af6b13"},
			wantErr: false,
		},
		{
			args:    []string{constants.DeleteCmd, constants.TagCmd, "-t", "invalid id"},
			wantErr: true,
		},
		{
			args:    []string{constants.DeleteCmd, constants.TagCmd, "-t", ""},
			wantErr: true,
		},
		{
			args:    []string{constants.DeleteCmd, constants.TagCmd, "-q", "@#$invalid-id", "-t", "23011406-6f3b-4431-9363-4e1af9af6b13"},
			wantErr: true,
		},
	}

	deleteCmd.AddCommand(deleteTagCmd)
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

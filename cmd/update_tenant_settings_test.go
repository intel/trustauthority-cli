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

func TestUpdateTenantSettingsCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)
	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.UpdateCmd, constants.TenantSettingsCmd},
			wantErr: true,
		},
		{
			args:    []string{constants.UpdateCmd, constants.TenantSettingsCmd, "-q", "valid-id", "-e", "dummy@email.com"},
			wantErr: false,
		},
		{
			args:    []string{constants.UpdateCmd, constants.TenantSettingsCmd, "-e", "invalid-email/script"},
			wantErr: true,
		},
		{
			args:    []string{constants.UpdateCmd, constants.TenantSettingsCmd, "-d"},
			wantErr: false,
		},
		{
			args:    []string{constants.UpdateCmd, constants.TenantSettingsCmd, "-q", "@#$invalid-id", "-e", "dummy@email.com"},
			wantErr: true,
		},
	}

	updateCmd.AddCommand(updateTenantSettingsCmd)
	tenantCmd.AddCommand(updateCmd)

	for _, tc := range tt {
		_, err := execute(t, tenantCmd, tc.args)

		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
	assert.NoError(t, err)
}

func TestUpdateTenantCommandWithInvalidUrl(t *testing.T) {
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
			args:        []string{constants.UpdateCmd, constants.TenantSettingsCmd, "-q", "valid-id", "-e", "dummy@email.com"},
			wantErr:     true,
			url:         "bogus\nbase\nURL",
			description: "Test update tenant settings using invalid URL",
		},
		{
			args:        []string{constants.UpdateCmd, constants.TenantSettingsCmd, "-q", "valid-id", "-e", "dummy@email.com"},
			wantErr:     true,
			url:         "a/b/c",
			description: "Invalid send request provided for update tenant settings command ",
		},
	}

	updateCmd.AddCommand(updateTenantSettingsCmd)
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

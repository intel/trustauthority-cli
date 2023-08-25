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

func TestDeletePolicyCommandWithInvalidUrl(t *testing.T) {
	test.SetupMockConfiguration("invalid url", tempConfigFile)
	load, err := config.LoadConfiguration()
	assert.NoError(t, err)
	viper.Set("trustauthority-url", "bogus\nbase\nURL")

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
	viper.Set("trustauthority-url", load.TrustAuthorityBaseUrl)
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
			args:    []string{constants.DeleteCmd, constants.PolicyCmd, "-q", "valid-id", "-p", "e48dabc5-9608-4ff3-aaed-f25909ab9de1"},
			wantErr: false,
		},
		{
			args:        []string{constants.DeleteCmd, constants.PolicyCmd, "-p", "invalid id"},
			wantErr:     true,
			description: "Test Invalid policy id provided",
		},
		{
			args:        []string{constants.DeleteCmd, constants.PolicyCmd, "-q", "@#$invalid-id", "-p", "e48dabc5-9608-4ff3-aaed-f25909ab9de1"},
			wantErr:     true,
			description: "Test Invalid request id provided",
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

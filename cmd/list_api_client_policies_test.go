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

func TestListApiClientCommandWithInvalidUrl(t *testing.T) {
	test.SetupMockConfiguration("invalid url", tempConfigFile)
	load, err := config.LoadConfiguration()
	assert.NoError(t, err)
	viper.Set("amber-base-url", "bogus\nbase\nURL")

	invalidUrlTc := struct {
		args        []string
		wantErr     bool
		description string
	}{
		args: []string{constants.ListCmd, constants.ApiClientCmd, constants.PolicyCmd, "-r",
			"5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-c", "3780cc39-cce2-4ec2-a47f-03e55b12e259"},
		wantErr:     true,
		description: "Test list api client policies using invalid URL",
	}
	getApiClientsCmd.AddCommand(getApiClientPoliciesCmd)
	listCmd.AddCommand(getApiClientsCmd)
	tenantCmd.AddCommand(listCmd)

	_, err = execute(t, tenantCmd, invalidUrlTc.args)
	viper.Set("amber-base-url", load.AmberBaseUrl)
	assert.Error(t, err)
}
func TestListApiClientsPoliciesCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args: []string{constants.ListCmd, constants.ApiClientCmd, constants.PolicyCmd, "-r",
				"5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-c", "3780cc39-cce2-4ec2-a47f-03e55b12e259"},
			wantErr: false,
		},
		{
			args: []string{constants.ListCmd, constants.ApiClientCmd, constants.PolicyCmd, "-r",
				"invalid id", "-c", "3780cc39-cce2-4ec2-a47f-03e55b12e259"},
			wantErr:     true,
			description: "Test Invalid service id provided",
		},
		{
			args: []string{constants.ListCmd, constants.ApiClientCmd, constants.PolicyCmd, "-r",
				"5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-c", "invalid id"},
			wantErr:     true,
			description: "Test Invalid apiClient id provided",
		},
	}

	getApiClientsCmd.AddCommand(getApiClientPoliciesCmd)
	listCmd.AddCommand(getApiClientsCmd)
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

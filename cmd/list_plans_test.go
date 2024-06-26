/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 *
 *
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

func TestListPlansCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.ListCmd, constants.PlanCmd, "-q", "valid-id", "-r", "ee28f3c2-6f58-489d-aa46-1140565d4718"},
			wantErr: false,
		},
		{
			args: []string{constants.ListCmd, constants.PlanCmd, "-r", "ee28f3c2-6f58-489d-aa46-1140565d4718", "-p",
				"a3ad72aa-86d6-49aa-b851-6a52caf0941b"},
			wantErr: false,
		},
		{
			args: []string{constants.ListCmd, constants.PlanCmd, "-r", "invalid id", "-p",
				"a3ad72aa-86d6-49aa-b851-6a52caf0941b"},
			wantErr:     true,
			description: "Invalid service offer id provided",
		},
		{
			args: []string{constants.ListCmd, constants.PlanCmd, "-r", "ee28f3c2-6f58-489d-aa46-1140565d4718", "-p",
				"invalid id"},
			wantErr:     true,
			description: "Invalid plan id provided",
		},
		{
			args:        []string{constants.ListCmd, constants.PlanCmd, "-q", "@#$valid-id", "-r", "ee28f3c2-6f58-489d-aa46-1140565d4718"},
			wantErr:     true,
			description: "Invalid request id provided",
		},
	}

	listCmd.AddCommand(getPlansCmd)
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

func TestListPlansCommandWithInvalidUrl(t *testing.T) {
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
			args:        []string{constants.ListCmd, constants.PlanCmd, "-q", "valid-id", "-r", "ee28f3c2-6f58-489d-aa46-1140565d4718"},
			wantErr:     true,
			url:         "bogus\nbase\nURL",
			description: "Test list plans using invalid URL",
		},
		{
			args:        []string{constants.ListCmd, constants.PlanCmd, "-q", "valid-id", "-r", "ee28f3c2-6f58-489d-aa46-1140565d4718"},
			wantErr:     true,
			url:         "a/b/c",
			description: "Invalid send request provided for list plans command ",
		},
	}

	listCmd.AddCommand(getPlansCmd)
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

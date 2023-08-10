/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"github.com/stretchr/testify/assert"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/test"
	"testing"
)

func TestListPoliciesCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:        []string{constants.ListCmd, constants.PolicyCmd, "-q", "valid-id"},
			wantErr:     false,
			description: "Get all policies under a tenant",
		},
		{
			args: []string{constants.ListCmd, constants.PolicyCmd, "-p",
				"e48dabc5-9608-4ff3-aaed-f25909ab9de1"},
			wantErr:     false,
			description: "Test Retrieve a policy under a tenant",
		},
		{
			args: []string{constants.ListCmd, constants.PolicyCmd, "-p",
				"invalid policy id"},
			wantErr:     true,
			description: "Test invalid policy id provided",
		},
		{
			args:        []string{constants.ListCmd, constants.PolicyCmd, "-q", "@#$invalid-id"},
			wantErr:     true,
			description: "Test invalid request id provided",
		},
	}

	listCmd.AddCommand(getPoliciesCmd)
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

/*
 * Copyright (C) 2023 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"github.com/stretchr/testify/assert"
	"intel/tac/v1/constants"
	"intel/tac/v1/test"
	"testing"
)

func TestListTenantSettingsCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)
	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.ListCmd, constants.TenantSettingsCmd, "-q", "valid-id"},
			wantErr: false,
		},
		{
			args:        []string{constants.ListCmd, constants.TenantSettingsCmd, "-q", "@#$invalid-id"},
			wantErr:     true,
			description: "Test invalid request id provided",
		},
	}

	listCmd.AddCommand(listTenantSettingsCmd)
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

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

func TestUpdateUserRoleCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args: []string{constants.UpdateCmd, constants.UserCmd, constants.RoleCmd, "-a", "abc",
				"-r", "User", "-u", "23011406-6f3b-4431-9363-4e1af9af6b13"},
			wantErr: false,
		},
	}

	updateUserCmd.AddCommand(updateUserRoleCmd)
	updateCmd.AddCommand(updateUserCmd)
	tenantCmd.AddCommand(updateCmd)

	for _, tc := range tt {
		_, err := execute(t, tenantCmd, tc.args)

		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

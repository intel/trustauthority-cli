/*
 * Copyright (C) 2023 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"github.com/stretchr/testify/assert"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/test"
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
			args:    []string{constants.UpdateCmd, constants.TenantSettingsCmd, "-e", "dummy@email.com"},
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

/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"intel/tac/v1/constants"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUninstallCmd(t *testing.T) {

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.UninstallCmd},
			wantErr: false,
		},
	}

	tenantCmd.AddCommand(uninstallCmd)

	for _, tc := range tt {
		_, err := execute(t, tenantCmd, tc.args)

		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestUninstallFunction(t *testing.T) {
	// Test the uninstall function directly
	err := uninstall()
	// Should not error even if directories don't exist
	assert.NoError(t, err)
}

func TestUninstallCmdWithErrors(t *testing.T) {
	// Test uninstall command execution
	tt := []struct {
		name        string
		args        []string
		wantErr     bool
		description string
	}{
		{
			name:        "Uninstall with valid command",
			args:        []string{constants.UninstallCmd},
			wantErr:     false,
			description: "Should complete successfully",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := execute(t, tenantCmd, tc.args)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

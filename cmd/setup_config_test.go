/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"github.com/stretchr/testify/assert"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/test"
	"os"
	"testing"
)

func TestConfigCmd(t *testing.T) {

	server := test.MockServer(t)
	os.Setenv("AMBER_BASE_URL", server.URL)
	os.Setenv("TENANT_ID", "f04971b7-fb41-4a9e-a06e-4bf6e71f98b3")

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.SetupConfigCmd, "-v", tempConfigFile.Name()},
			wantErr: false,
		},
		{
			args:    []string{constants.SetupConfigCmd, "-v"},
			wantErr: true,
		},
	}

	tenantCmd.AddCommand(setupConfigCmd)

	for _, tc := range tt {
		_, err := execute(t, tenantCmd, tc.args)

		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

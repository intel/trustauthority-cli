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

func TestListApiClientsCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.ListCmd, constants.ApiClientCmd, "-a", "abc", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777"},
			wantErr: false,
		},
		{
			args: []string{constants.ListCmd, constants.ApiClientCmd, "-a", "abc", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-c",
				"3780cc39-cce2-4ec2-a47f-03e55b12e259"},
			wantErr: false,
		},
	}

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

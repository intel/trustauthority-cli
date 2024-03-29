/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"github.com/stretchr/testify/assert"
	"intel/tac/v1/constants"
	"intel/tac/v1/test"
	"testing"
)

func TestListServicesCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.ListCmd, constants.ServiceCmd, "-q", "valid-id"},
			wantErr: false,
		},
		{
			args:    []string{constants.ListCmd, constants.ServiceCmd, "-r", "ae3d7720-08ab-421c-b8d4-1725c358f03e"},
			wantErr: false,
		},
		{
			args:        []string{constants.ListCmd, constants.ServiceCmd, "-r", "invalid id"},
			wantErr:     true,
			description: "Invalid service id provided",
		},
		{
			args:        []string{constants.ListCmd, constants.ServiceCmd, "-q", "@#$invalid-id"},
			wantErr:     true,
			description: "Invalid request id provided",
		},
	}

	listCmd.AddCommand(getServicesCmd)
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

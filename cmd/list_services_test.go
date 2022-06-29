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

func TestListServicesCmd(t *testing.T) {
	server := test.MockTmsServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.ListCmd, constants.ServiceCmd, "-a", "abc"},
			wantErr: false,
		},
		{
			args:    []string{constants.ListCmd, constants.ServiceCmd, "-a", "abc", "-r", "ae3d7720-08ab-421c-b8d4-1725c358f03e"},
			wantErr: false,
		},
	}

	listCmd.AddCommand(getServicesCmd)
	tenantCmd.AddCommand(listCmd)

	for _, tc := range tt {
		_, err := execute(t, tenantCmd, tc.args...)

		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

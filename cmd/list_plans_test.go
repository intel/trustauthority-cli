/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 *
 *
 */

package cmd

import (
	"github.com/stretchr/testify/assert"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/test"
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
			args:    []string{constants.ListCmd, constants.PlanCmd, "-r", "ee28f3c2-6f58-489d-aa46-1140565d4718"},
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

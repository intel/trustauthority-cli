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

func TestUpdatePolicyCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)
	GenerateInvalidPolicyFile(t, tempPolicyFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args: []string{constants.UpdateCmd, constants.PolicyCmd, "-i", "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
				"-n", "Sample_Policy_SGX", "-f", "../test/resources/rego-policy.txt"},
			wantErr: false,
		},
		{

			args: []string{constants.UpdateCmd, constants.PolicyCmd, "-i", "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
				"-n", "Sample_Policy_SGX", "-f", tempPolicyFile},
			wantErr:     true,
			description: "Test Create Policy With Invalid File Size",
		},
		{
			args: []string{constants.UpdateCmd, constants.PolicyCmd, "-i", "invalid id",
				"-n", "Sample_Policy_SGX", "-f", "../test/resources/rego-policy.txt"},
			wantErr:     true,
			description: "Test Invalid policy Id provided",
		},
		{
			args: []string{constants.UpdateCmd, constants.PolicyCmd, "-i", "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
				"-n", "Sample_Policy_SGX", "-f", ""},
			wantErr:     true,
			description: "Test Policy file path cannot be empty",
		},
		{
			args: []string{constants.UpdateCmd, constants.PolicyCmd, "-i", "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
				"-n", "Sample_Policy_SGX", "-f", "../test/resources/@rego-policy.txt"},
			wantErr:     true,
			description: "Test Unsafe or invalid path specified",
		},
	}

	updateCmd.AddCommand(updatePolicyCmd)
	tenantCmd.AddCommand(updateCmd)

	for _, tc := range tt {
		_, err := execute(t, tenantCmd, tc.args)

		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
	err = os.Remove(tempPolicyFile)
	assert.NoError(t, err)
}

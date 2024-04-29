/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	"intel/tac/v1/test"
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
			args: []string{constants.UpdateCmd, constants.PolicyCmd, "-q", "valid-id", "-i", "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
				"-n", "Sample_Policy_SGX", "-f", "../test/resources/rego-policy.txt"},
			wantErr: false,
		},
		{
			args: []string{constants.CreateCmd, constants.PolicyCmd, "-n", "Sample Policy SGX", "-t", "Appraisal policy",
				"-r", "e8a72b7e-c4b1-4bdc-bf40-68f23c68a2aa", "-a", "SGX Attestation", "-f", "../test/resources/rego-policy.txt"},
			wantErr:     true,
			description: "Test Create Policy with invalid policy name",
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
			wantErr:     false,
			description: "Test Policy file empty",
		},
		{
			args: []string{constants.UpdateCmd, constants.PolicyCmd, "-i", "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
				"-n", "Sample_Policy_SGX", "-f", "../test/resources/@rego-policy.txt"},
			wantErr:     true,
			description: "Test Unsafe or invalid path specified",
		},
		{
			args: []string{constants.UpdateCmd, constants.PolicyCmd, "-i", "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
				"-n", "Sample_Policy_SGX", "-f", "../test/resources/"},
			wantErr:     true,
			description: "Test Error reading policy file",
		},
		{
			args: []string{constants.UpdateCmd, constants.PolicyCmd, "-q", "@#$invalid-id", "-i", "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
				"-n", "Sample_Policy_SGX", "-f", "../test/resources/rego-policy.txt"},
			wantErr:     true,
			description: "Test invalid request id provided",
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

func TestUpdatePolicyCommandWithInvalidUrl(t *testing.T) {
	test.SetupMockConfiguration("invalid url", tempConfigFile)
	load, err := config.LoadConfiguration()
	assert.NoError(t, err)

	invalidUrlTc := []struct {
		args        []string
		wantErr     bool
		url         string
		description string
	}{
		{
			args: []string{constants.UpdateCmd, constants.PolicyCmd, "-q", "valid-id", "-i", "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
				"-n", "Sample_Policy_SGX", "-f", "../test/resources/rego-policy.txt"},
			wantErr:     true,
			url:         "bogus\nbase\nURL",
			description: "Test update policy using invalid URL",
		},
		{
			args: []string{constants.UpdateCmd, constants.PolicyCmd, "-q", "valid-id", "-i", "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
				"-n", "Sample_Policy_SGX", "-f", "../test/resources/rego-policy.txt"},
			wantErr:     true,
			url:         "a/b/c",
			description: "Invalid send request provided for update policy command ",
		},
	}

	updateCmd.AddCommand(updatePolicyCmd)
	tenantCmd.AddCommand(updateCmd)

	for _, tc := range invalidUrlTc {
		viper.Set("trustauthority-url", tc.url)
		_, err := execute(t, tenantCmd, tc.args)
		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
	viper.Set("trustauthority-url", load.TrustAuthorityBaseUrl)
}

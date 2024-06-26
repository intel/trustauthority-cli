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
	"testing"
)

func TestUpdateApiClientCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args: []string{constants.UpdateCmd, constants.ApiClientCmd, "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i",
				"5f7eece7-ab3f-4f1f-98cd-31c6a44a9900,116690e2-ddf7-45b9-b943-35744cd34717", "-v",
				"d021545b-5f50-4923-8ad7-83274b3761b0:tag1,60fbbfbf-cd70-4c88-bedb-c057e3dca626:tag2",
				"-c", "3780cc39-cce2-4ec2-a47f-03e55b12e259", "-s", "Invalid_Status"},
			wantErr:     true,
			description: "Test Invalid status. should be one of Active, Inactive or Cancelled",
		},
		{
			args: []string{constants.UpdateCmd, constants.ApiClientCmd, "-q", "valid-id", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i",
				"5f7eece7-ab3f-4f1f-98cd-31c6a44a9900,116690e2-ddf7-45b9-b943-35744cd34717", "-v",
				"d021545b-5f50-4923-8ad7-83274b3761b0:tag1,60fbbfbf-cd70-4c88-bedb-c057e3dca626:tag2",
				"-c", "3780cc39-cce2-4ec2-a47f-03e55b12e259", "-s", "Active"},
			wantErr: false,
		},
		{
			args: []string{constants.UpdateCmd, constants.ApiClientCmd, "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
				"-v", "5cfb6af4-59ac-4a14-8b83-bd65b1e11779"},
			wantErr:     true,
			description: "Test Create ApiClient With Incorrect Tag Format",
		},
		{
			args: []string{constants.CreateCmd, constants.ApiClientCmd, "-a", "abc", "-n", "Test_Subs", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-v", "invalid Tag Name:tag1"},
			wantErr:     true,
			description: "Test Create Api Client With Incorrect Tag Name",
		},
		{
			args: []string{constants.CreateCmd, constants.ApiClientCmd, "-a", "abc", "-n", "Test_Subs", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-v", "tagName:invalid value"},
			wantErr:     true,
			description: "Test Create Api Client With Incorrect Tag Value",
		},
		{
			args: []string{constants.UpdateCmd, constants.ApiClientCmd, "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i", "abc"},
			wantErr:     true,
			description: "Test Create ApiClient With Invalid Policy IDs",
		},
		{
			args: []string{constants.UpdateCmd, constants.ApiClientCmd, "-p",
				"invalid id ", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i",
				"5f7eece7-ab3f-4f1f-98cd-31c6a44a9900,116690e2-ddf7-45b9-b943-35744cd34717", "-v",
				"d021545b-5f50-4923-8ad7-83274b3761b0:tag1,60fbbfbf-cd70-4c88-bedb-c057e3dca626:tag2",
				"-c", "3780cc39-cce2-4ec2-a47f-03e55b12e259", "-s", "Active"},
			wantErr:     true,
			description: "Test Invalid product Id provided",
		},
		{
			args: []string{constants.UpdateCmd, constants.ApiClientCmd, "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "invalid id", "-i",
				"5f7eece7-ab3f-4f1f-98cd-31c6a44a9900,116690e2-ddf7-45b9-b943-35744cd34717", "-v",
				"d021545b-5f50-4923-8ad7-83274b3761b0:tag1,60fbbfbf-cd70-4c88-bedb-c057e3dca626:tag2",
				"-c", "3780cc39-cce2-4ec2-a47f-03e55b12e259", "-s", "Active"},
			wantErr:     true,
			description: "Test Invalid service Id provided",
		},
		{
			args: []string{constants.UpdateCmd, constants.ApiClientCmd, "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i",
				"5f7eece7-ab3f-4f1f-98cd-31c6a44a9900,116690e2-ddf7-45b9-b943-35744cd34717", "-v",
				"d021545b-5f50-4923-8ad7-83274b3761b0:tag1,60fbbfbf-cd70-4c88-bedb-c057e3dca626:tag2",
				"-c", "invalid id", "-s", "Active"},
			wantErr:     true,
			description: "Test Invalid api client Id provided",
		},
		{
			args: []string{constants.UpdateCmd, constants.ApiClientCmd, "-q", "@#$invalid-id", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i",
				"5f7eece7-ab3f-4f1f-98cd-31c6a44a9900,116690e2-ddf7-45b9-b943-35744cd34717", "-v",
				"d021545b-5f50-4923-8ad7-83274b3761b0:tag1,60fbbfbf-cd70-4c88-bedb-c057e3dca626:tag2",
				"-c", "3780cc39-cce2-4ec2-a47f-03e55b12e259", "-s", "Active"},
			wantErr: true,
		},
	}

	updateCmd.AddCommand(updateApiClientCmd)
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

func TestUpdateAPIClieniWithInvalidUrl(t *testing.T) {
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
			args: []string{constants.UpdateCmd, constants.ApiClientCmd, "-q", "valid-id", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i",
				"5f7eece7-ab3f-4f1f-98cd-31c6a44a9900,116690e2-ddf7-45b9-b943-35744cd34717", "-v",
				"d021545b-5f50-4923-8ad7-83274b3761b0:tag1,60fbbfbf-cd70-4c88-bedb-c057e3dca626:tag2",
				"-c", "3780cc39-cce2-4ec2-a47f-03e55b12e259", "-s", "Active"},
			url:         "bogus\nbase\nURL",
			wantErr:     true,
			description: "Update API Client with Invalid URL",
		},
		{
			args: []string{constants.UpdateCmd, constants.ApiClientCmd, "-q", "valid-id", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i",
				"5f7eece7-ab3f-4f1f-98cd-31c6a44a9900,116690e2-ddf7-45b9-b943-35744cd34717", "-v",
				"d021545b-5f50-4923-8ad7-83274b3761b0:tag1,60fbbfbf-cd70-4c88-bedb-c057e3dca626:tag2",
				"-c", "3780cc39-cce2-4ec2-a47f-03e55b12e259", "-s", "Active"},
			wantErr:     true,
			url:         "a/b/c",
			description: "Invalid send request provided for update api client command",
		},
	}

	createCmd.AddCommand(updateApiClientCmd)
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

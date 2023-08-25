/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	"intel/tac/v1/test"
	"testing"
)

func TestCreateApiClientCommandWithInvalidUrl(t *testing.T) {
	server_r := test.MockServer(t)
	test.SetupMockConfiguration(server_r.URL, tempConfigFile)
	load, err := config.LoadConfiguration()
	assert.NoError(t, err)
	viper.Set("trustauthority-url", "bogus\nbase\nURL")

	invalidUrlTc := struct {
		args        []string
		wantErr     bool
		description string
	}{
		args: []string{constants.CreateCmd, constants.ApiClientCmd, "-n", "Test_Subs", "-p",
			"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i", "5f7eece7-ab3f-4f1f-98cd-31c6a44a9900",
			"-v", "Workload:WorkloadAI,Workload:WorkloadEXE"},
		wantErr:     true,
		description: "Test Create Api Client using invalid URL",
	}

	createCmd.AddCommand(createApiClientCmd)
	tenantCmd.AddCommand(createCmd)

	_, err = execute(t, tenantCmd, invalidUrlTc.args)
	viper.Set("trustauthority-url", load.TrustAuthorityBaseUrl)
	assert.Error(t, err)
}

func TestCreateApiClientCmd(t *testing.T) {
	server := test.MockServer(t)
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args: []string{constants.CreateCmd, constants.ApiClientCmd, "-q", "valid-id", "-n", "Test_Subs", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i", "5f7eece7-ab3f-4f1f-98cd-31c6a44a9900",
				"-v", "Workload:WorkloadAI,Workload:WorkloadEXE"},
			wantErr:     false,
			description: "Test Create Api Client",
		},
		{
			args: []string{constants.CreateCmd, constants.ApiClientCmd, "-n", "Test_Subs", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
				"-p", "invalid id"},
			wantErr:     true,
			description: "Test Invalid product id provided",
		},
		{
			args: []string{constants.CreateCmd, constants.ApiClientCmd, "-n", "Test_Subs", "-r", "invalid id",
				"-p", "e169d34f-58ce-4717-9b3a-5c66abd33417"},
			wantErr:     true,
			description: "Test Invalid service id provided",
		},
		{
			args: []string{constants.CreateCmd, constants.ApiClientCmd, "-n", "@@@@@@@", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
				"-p", "e169d34f-58ce-4717-9b3a-5c66abd33417"},
			wantErr:     true,
			description: "Test Invalid api client name provided",
		},
		{
			args: []string{constants.CreateCmd, constants.ApiClientCmd, "-n", "Test_Subs", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777"},
			wantErr:     false,
			description: "Test Create Api Client",
		},
		{
			args: []string{constants.CreateCmd, constants.ApiClientCmd, "-n", "Test_Subs", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
				"-v", "tagName"},
			wantErr:     true,
			description: "Test Create Api Client With Incorrect Tag Format",
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
			args: []string{constants.CreateCmd, constants.ApiClientCmd, "-n", "Test_Subs", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i", "invalid id"},
			wantErr:     true,
			description: "Test Create Api Client With Invalid Policy IDs",
		},
		{
			args: []string{constants.CreateCmd, constants.ApiClientCmd, "-q", "@#$invalid-id", "Test_Subs", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i", "5f7eece7-ab3f-4f1f-98cd-31c6a44a9900",
				"-v", "Workload:WorkloadAI,Workload:WorkloadEXE"},
			wantErr:     true,
			description: "Test Create Api Client with invalid request ID",
		},
	}

	createCmd.AddCommand(createApiClientCmd)
	tenantCmd.AddCommand(createCmd)

	for _, tc := range tt {
		fmt.Println(tc.args)
		_, err := execute(t, tenantCmd, tc.args)

		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

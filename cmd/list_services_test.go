/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"intel/tac/v1/config"
	"intel/tac/v1/constants"

	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestListServicesCmd(t *testing.T) {
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

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

func TestListServicesTestCommandWithInvalidUrl(t *testing.T) {
	setupMockConfiguration("invalid url", tempConfigFile)
	load, err := config.LoadConfiguration()
	assert.NoError(t, err)

	invalidUrlTc := []struct {
		args        []string
		wantErr     bool
		url         string
		description string
	}{
		{
			args:        []string{constants.ListCmd, constants.ServiceCmd, "-q", "valid-id"},
			wantErr:     true,
			url:         "bogus\\nbase\\nURL",
			description: "Test list services using invalid URL",
		},
		{
			args:        []string{constants.ListCmd, constants.ServiceCmd, "-q", "valid-id"},
			wantErr:     true,
			url:         "a/b/c",
			description: "Invalid send request provided for list services command ",
		},
	}

	listCmd.AddCommand(getServicesCmd)
	tenantCmd.AddCommand(listCmd)

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

func TestListServicesWithServiceId(t *testing.T) {
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:        []string{constants.ListCmd, constants.ServiceCmd, "-r", "ae3d7720-08ab-421c-b8d4-1725c358f03e"},
			wantErr:     false,
			description: "List specific service by ID",
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

func TestListServicesAllCases(t *testing.T) {
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		name        string
		args        []string
		wantErr     bool
		description string
	}{
		{
			name:        "List all services",
			args:        []string{constants.ListCmd, constants.ServiceCmd},
			wantErr:     false,
			description: "Should list all services when no ID provided",
		},
		{
			name:        "List service with valid UUID",
			args:        []string{constants.ListCmd, constants.ServiceCmd, "-r", "ae3d7720-08ab-421c-b8d4-1725c358f03e"},
			wantErr:     false,
			description: "Should retrieve specific service",
		},
		{
			name:        "List service with valid request ID",
			args:        []string{constants.ListCmd, constants.ServiceCmd, "-q", "test-request-id"},
			wantErr:     false,
			description: "Should work with custom request ID",
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

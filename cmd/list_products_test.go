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

func TestListProductsCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.ListCmd, constants.ProductCmd, "-q", "valid-id", "-r", "ae3d7720-08ab-421c-b8d4-1725c358f03e"},
			wantErr: false,
		},
		{
			args:        []string{constants.ListCmd, constants.ProductCmd, "-r", "invalid id"},
			wantErr:     true,
			description: "Test Invalid Service offer Id provided",
		},
		{
			args:        []string{constants.ListCmd, constants.ProductCmd, "-q", "@#$invalid-id", "-r", "ae3d7720-08ab-421c-b8d4-1725c358f03e"},
			wantErr:     true,
			description: "Test Invalid Request Id provided",
		},
	}

	listCmd.AddCommand(getProductsCmd)
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

func TestListApiProductCommnadWithInvalidUrl(t *testing.T) {
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
			args:        []string{constants.ListCmd, constants.ProductCmd, "-q", "valid-id", "-r", "ae3d7720-08ab-421c-b8d4-1725c358f03e"},
			wantErr:     true,
			url:         "bogus\nbase\nURL",
			description: "Test list products using invalid URL",
		},
		{
			args:        []string{constants.ListCmd, constants.ProductCmd, "-q", "valid-id", "-r", "ae3d7720-08ab-421c-b8d4-1725c358f03e"},
			wantErr:     true,
			url:         "a/b/c",
			description: "Invalid send request provided for list products command ",
		},
	}

	listCmd.AddCommand(getProductsCmd)
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

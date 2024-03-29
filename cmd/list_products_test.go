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

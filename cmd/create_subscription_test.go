/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/test"
	"testing"
)

func TestCreateSubscriptionCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args: []string{constants.CreateCmd, constants.SubscriptionCmd, "-a", "abc", "-n", "Test Subs", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-e", "2022-09-24"},
			wantErr:     false,
			description: "Test Create Subscription",
		},
		{
			args: []string{constants.CreateCmd, constants.SubscriptionCmd, "-a", "abc", "-n", "Test Subs", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
				"-v", "5cfb6af4-59ac-4a14-8b83-bd65b1e11779"},
			wantErr:     true,
			description: "Test Create Subscription With Incorrect Tag Format",
		},
		{
			args: []string{constants.CreateCmd, constants.SubscriptionCmd, "-a", "abc", "-n", "Test Subs", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-v", "invalid-UUID:tag1"},
			wantErr:     true,
			description: "Test Create Subscription With Incorrect Tag UUID",
		},
		{
			args: []string{constants.CreateCmd, constants.SubscriptionCmd, "-a", "abc", "-d", "Test Subs", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i", "abc"},
			wantErr:     true,
			description: "Test Create Subscription With Invalid Policy IDs",
		},
	}

	createCmd.AddCommand(createSubscriptionCmd)
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

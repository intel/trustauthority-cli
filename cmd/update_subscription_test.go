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

func TestUpdateSubscriptionCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args: []string{constants.UpdateCmd, constants.SubscriptionCmd, "-a", "abc", "-n", "Test Subs", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i",
				"5f7eece7-ab3f-4f1f-98cd-31c6a44a9900,116690e2-ddf7-45b9-b943-35744cd34717", "-v",
				"d021545b-5f50-4923-8ad7-83274b3761b0:tag1,60fbbfbf-cd70-4c88-bedb-c057e3dca626:tag2",
				"-u", "3780cc39-cce2-4ec2-a47f-03e55b12e259", "-s", "Active", "-e", "2022-09-25"},
			wantErr: false,
		},
		{
			args: []string{constants.UpdateCmd, constants.SubscriptionCmd, "-a", "abc", "-n", "Test Subs", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
				"-v", "5cfb6af4-59ac-4a14-8b83-bd65b1e11779"},
			wantErr:     true,
			description: "Test Create Subscription With Incorrect Tag Format",
		},
		{
			args: []string{constants.UpdateCmd, constants.SubscriptionCmd, "-a", "abc", "-n", "Test Subs", "-p",
				"e169d34f-58ce-4717-9b3a-5c66abd33417", "-r", "5cfb6af4-59ac-4a14-8b83-bd65b1e11777", "-i", "abc"},
			wantErr:     true,
			description: "Test Create Subscription With Invalid Policy IDs",
		},
	}

	updateCmd.AddCommand(updateSubscriptionCmd)
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

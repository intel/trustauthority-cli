/*
 * Copyright (C) 2026 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDeleteRimWithNoParams(t *testing.T) {
	// Arrange
	deleteCmd.AddCommand(deleteRimCmd)
	tenantCmd.AddCommand(deleteCmd)
	args := []string{constants.DeleteCmd, constants.RimCmd}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "required flag(s) \"rim-id\" not set")
}

func TestDeleteRimWithInvalidUUID(t *testing.T) {
	// Arrange
	deleteCmd.AddCommand(deleteRimCmd)
	tenantCmd.AddCommand(deleteCmd)
	args := []string{constants.DeleteCmd, constants.RimCmd, "-r", "invalid-uuid"}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "Invalid RIM id provided: invalid UUID length: 12")
}

func TestDeleteRimWithValidUUID(t *testing.T) {
	// Arrange
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	validUUID := "e1e4424b-85cc-41bb-b295-7a24c3e8a8aa"
	deleteCmd.AddCommand(deleteRimCmd)
	tenantCmd.AddCommand(deleteCmd)
	args := []string{constants.DeleteCmd, constants.RimCmd, "-r", validUUID}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err)
}

func TestDeleteRimWithInvalidUrl(t *testing.T) {
	// Arrange
	validUUID := uuid.New().String()
	load, err := config.LoadConfiguration()
	if err == nil {
		viper.Set("trustauthority-url", "bogus\nbase\nURL")
		defer viper.Set("trustauthority-url", load.TrustAuthorityBaseUrl)

		deleteCmd.AddCommand(deleteRimCmd)
		tenantCmd.AddCommand(deleteCmd)
		args := []string{constants.DeleteCmd, constants.RimCmd, "-r", validUUID}

		// Act
		_, err = execute(t, tenantCmd, args)

		// Assert
		assert.ErrorContains(t, err, "url")
	}
}

func TestDeleteRimWithEmptyUUID(t *testing.T) {
	// Arrange
	deleteCmd.AddCommand(deleteRimCmd)
	tenantCmd.AddCommand(deleteCmd)
	args := []string{constants.DeleteCmd, constants.RimCmd, "-r", ""}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "Invalid RIM id provided: invalid UUID length: 0")
}

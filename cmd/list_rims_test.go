/*
 * Copyright (C) 2026 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"intel/tac/v1/constants"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListRimsWithValidUUID(t *testing.T) {
	// Arrange
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	validUUID := "a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6" // Matches mock RIM data
	listCmd.AddCommand(getRimsCmd)
	tenantCmd.AddCommand(listCmd)
	args := []string{constants.ListCmd, constants.RimCmd, "-r", validUUID}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err)
}

func TestListRimsWithNameFilter(t *testing.T) {
	// Arrange
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	listCmd.AddCommand(getRimsCmd)
	tenantCmd.AddCommand(listCmd)
	args := []string{constants.ListCmd, constants.RimCmd, "-n", "test.rim"}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err)
}

func TestListAllRims(t *testing.T) {
	// Arrange
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	listCmd.AddCommand(getRimsCmd)
	tenantCmd.AddCommand(listCmd)
	args := []string{constants.ListCmd, constants.RimCmd}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err)
}

func TestListRimsWithIncludeOwnPrivate(t *testing.T) {
	// Arrange
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	listCmd.AddCommand(getRimsCmd)
	tenantCmd.AddCommand(listCmd)
	args := []string{constants.ListCmd, constants.RimCmd, constants.IncludeOwnPrivateParamName}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err)
}

func TestListRimsWithIncludeOwnPublic(t *testing.T) {
	// Arrange
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	listCmd.AddCommand(getRimsCmd)
	tenantCmd.AddCommand(listCmd)
	args := []string{constants.ListCmd, constants.RimCmd, constants.IncludeOwnPublicParamName}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err)
}

func TestListRimsWithIncludeOthersPublic(t *testing.T) {
	// Arrange
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	listCmd.AddCommand(getRimsCmd)
	tenantCmd.AddCommand(listCmd)
	args := []string{constants.ListCmd, constants.RimCmd, constants.IncludeOthersPublicParamName}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err)
}

func TestListRimsWithExcludeContent(t *testing.T) {
	// Arrange
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	listCmd.AddCommand(getRimsCmd)
	tenantCmd.AddCommand(listCmd)
	args := []string{constants.ListCmd, constants.RimCmd, constants.ExcludeContentParamName}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err)
}

func TestListRimsWithExcludeDescription(t *testing.T) {
	// Arrange
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	listCmd.AddCommand(getRimsCmd)
	tenantCmd.AddCommand(listCmd)
	args := []string{constants.ListCmd, constants.RimCmd, constants.ExcludeDescriptionParamName}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err)
}

func TestListRimsWithMultipleParams(t *testing.T) {
	// Arrange
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	listCmd.AddCommand(getRimsCmd)
	tenantCmd.AddCommand(listCmd)
	args := []string{constants.ListCmd, constants.RimCmd, "-n", "test.rim",
		constants.IncludeOwnPrivateParamName, constants.ExcludeContentParamName}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err)
}

func TestListRimsWithInvalidUUID(t *testing.T) {
	// Arrange
	listCmd.AddCommand(getRimsCmd)
	tenantCmd.AddCommand(listCmd)
	args := []string{constants.ListCmd, constants.RimCmd, "-r", "invalid-uuid"}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "Invalid RIM id provided: invalid UUID length: 12")
}

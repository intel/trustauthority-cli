/*
 * Copyright (C) 2026 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"intel/tac/v1/constants"
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

// setupUpdateRimTest prepares the command hierarchy and returns cleanup function
func setupUpdateRimTest(t *testing.T) {
	t.Helper()
	// Remove existing commands to ensure clean state
	tenantCmd.RemoveCommand(updateCmd)
	updateCmd.RemoveCommand(updateRimCmd)

	// Re-add commands
	updateCmd.AddCommand(updateRimCmd)
	tenantCmd.AddCommand(updateCmd)

	// Reset all flags
	updateRimCmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Value.Set(f.DefValue)
		f.Changed = false
	})

	// Cleanup function
	t.Cleanup(func() {
		updateRimCmd.Flags().VisitAll(func(f *pflag.Flag) {
			f.Value.Set(f.DefValue)
			f.Changed = false
		})
	})
}

func TestUpdateRimWithValidFile(t *testing.T) {
	// Arrange
	setupUpdateRimTest(t)
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	rimFilePath := "../test/resources/rim-policy.txt"
	validUUID := "e1e4424b-85cc-41bb-b295-7a24c3e8a8aa"
	args := []string{constants.UpdateCmd, constants.RimCmd, "-r", validUUID, "-f", rimFilePath}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err)
}

func TestUpdateRimWithNoParams(t *testing.T) {
	// Arrange
	setupUpdateRimTest(t)
	args := []string{constants.UpdateCmd, constants.RimCmd}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "required flag(s) \"rim-content-file\", \"rim-id\" not set")
}

func TestUpdateRimWithInvalidUUID(t *testing.T) {
	// Arrange
	setupUpdateRimTest(t)
	rimFilePath := "../test/resources/rim-policy.txt"
	args := []string{constants.UpdateCmd, constants.RimCmd, "-r", "invalid-uuid", "-f", rimFilePath}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "Invalid RIM id provided: invalid UUID length: 12")
}

func TestUpdateRimWithNoContentFile(t *testing.T) {
	// Arrange
	setupUpdateRimTest(t)
	args := []string{constants.UpdateCmd, constants.RimCmd, "-r", "e1e4424b-85cc-41bb-b295-7a24c3e8a8aa"}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "required flag(s) \"rim-content-file\" not set")
}

func TestUpdateRimWithInvalidFilePath(t *testing.T) {
	// Arrange
	setupUpdateRimTest(t)
	args := []string{constants.UpdateCmd, constants.RimCmd, "-r", "e1e4424b-85cc-41bb-b295-7a24c3e8a8aa", "-f", "/nonexistent/path/file.json"}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "Unsafe symlink detected in path")
}

func TestUpdateRimWithInvalidJSON(t *testing.T) {
	// Arrange
	setupUpdateRimTest(t)
	tmpFile, err := os.CreateTemp("", "invalid-rim-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("{invalid json")
	assert.NoError(t, err)
	tmpFile.Close()

	args := []string{constants.UpdateCmd, constants.RimCmd, "-r", "e1e4424b-85cc-41bb-b295-7a24c3e8a8aa", "-f", tmpFile.Name()}

	// Act
	_, err = execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "RIM content must be either valid JSON or valid JWT format")
}

func TestUpdateRimWithSignedJWT(t *testing.T) {
	// Arrange
	setupUpdateRimTest(t)
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	// Create a temp file with signed JWT content
	tmpFile, err := os.CreateTemp("", "rim-jwt-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Valid signed JWT with RIM payload
	jwtContent := "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZWFzdXJlbWVudHMiOlt7ImluZGV4IjowLCJ2YWx1ZSI6ImQ0MTJhNGYwN2VmODM4OTJhNTkxNWZiMmFiNTg0YmUzMWUxODZlNWE0Zjk1YWI1ZjY5NTBmZDRlYjg2OTRkN2IifSx7ImluZGV4IjoxLCJ2YWx1ZSI6ImJhYjkxZjIwMDAzODA3NmFjMjVmODdkZTBjYTY3NDcyNDQzYzJlYmUxN2VkOWJhOTUzMTRlNjA5MDM4ZjUxYWIifV0sIm1ldGFkYXRhIjp7InZlcnNpb24iOiIxLjAiLCJ0eXBlIjoicmVmZXJlbmNlIn19.signature"
	_, err = tmpFile.WriteString(jwtContent)
	assert.NoError(t, err)
	tmpFile.Close()

	validUUID := "e1e4424b-85cc-41bb-b295-7a24c3e8a8aa"
	args := []string{constants.UpdateCmd, constants.RimCmd, "-r", validUUID, "-f", tmpFile.Name()}

	// Act
	_, err = execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err, "Should accept signed JWT content for update")
}

func TestUpdateRimWithUnsignedJWT(t *testing.T) {
	// Arrange
	setupUpdateRimTest(t)
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	// Create a temp file with unsigned JWT content (alg=none)
	tmpFile, err := os.CreateTemp("", "rim-unsigned-jwt-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Valid unsigned JWT with RIM payload
	unsignedJWT := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJtZWFzdXJlbWVudHMiOlt7ImluZGV4IjowLCJ2YWx1ZSI6ImFiYzEyMyJ9XSwibWV0YWRhdGEiOnsidmVyc2lvbiI6IjEuMCJ9fQ."
	_, err = tmpFile.WriteString(unsignedJWT)
	assert.NoError(t, err)
	tmpFile.Close()

	validUUID := "e1e4424b-85cc-41bb-b295-7a24c3e8a8aa"
	args := []string{constants.UpdateCmd, constants.RimCmd, "-r", validUUID, "-f", tmpFile.Name()}

	// Act
	_, err = execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err, "Should accept unsigned JWT content (alg=none) for update")
}

func TestUpdateRimWithJWTAndWhitespace(t *testing.T) {
	// Arrange
	setupUpdateRimTest(t)
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	// Create a temp file with JWT and trailing whitespace
	tmpFile, err := os.CreateTemp("", "rim-jwt-whitespace-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// JWT with trailing newline and spaces
	jwtWithWhitespace := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZWFzdXJlbWVudHMiOltdfQ.sig  \n\n"
	_, err = tmpFile.WriteString(jwtWithWhitespace)
	assert.NoError(t, err)
	tmpFile.Close()

	validUUID := "e1e4424b-85cc-41bb-b295-7a24c3e8a8aa"
	args := []string{constants.UpdateCmd, constants.RimCmd, "-r", validUUID, "-f", tmpFile.Name()}

	// Act
	_, err = execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err, "Should handle JWT with whitespace for update")
}

func TestUpdateRimWithInvalidJWT(t *testing.T) {
	// Arrange
	setupUpdateRimTest(t)
	tmpFile, err := os.CreateTemp("", "invalid-jwt-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Invalid JWT - only 2 parts but malformed
	invalidJWT := "not-base64.also-not-base64"
	_, err = tmpFile.WriteString(invalidJWT)
	assert.NoError(t, err)
	tmpFile.Close()

	validUUID := "e1e4424b-85cc-41bb-b295-7a24c3e8a8aa"
	args := []string{constants.UpdateCmd, constants.RimCmd, "-r", validUUID, "-f", tmpFile.Name()}

	// Act
	_, err = execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "RIM content is not valid JSON or JWT format")
}

func TestUpdateRimSwitchingFromJSONToJWT(t *testing.T) {
	// Arrange
	setupUpdateRimTest(t)
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	// Create a temp file with JWT content (simulating updating from JSON to JWT)
	tmpFile, err := os.CreateTemp("", "rim-update-to-jwt-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// JWT content
	jwtContent := "eyJhbGciOiJSUzM4NCIsInR5cCI6IkpXVCJ9.eyJtZWFzdXJlbWVudHMiOlt7ImluZGV4IjowfV19.sig"
	_, err = tmpFile.WriteString(jwtContent)
	assert.NoError(t, err)
	tmpFile.Close()

	validUUID := "e1e4424b-85cc-41bb-b295-7a24c3e8a8aa"
	args := []string{constants.UpdateCmd, constants.RimCmd, "-r", validUUID, "-f", tmpFile.Name()}

	// Act
	_, err = execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err, "Should allow updating from plain JSON to JWT format")
}

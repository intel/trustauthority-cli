/*
 * Copyright (C) 2026 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// setupCreateRimTest prepares the command hierarchy and resets flag state
func setupCreateRimTest(t *testing.T) {
	t.Helper()
	// Remove existing commands to ensure clean state
	tenantCmd.RemoveCommand(createCmd)
	createCmd.RemoveCommand(createRimCmd)

	// Re-add commands
	createCmd.AddCommand(createRimCmd)
	tenantCmd.AddCommand(createCmd)

	// Reset all flags
	createRimCmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Value.Set(f.DefValue)
		f.Changed = false
	})

	// Cleanup function
	t.Cleanup(func() {
		createRimCmd.Flags().VisitAll(func(f *pflag.Flag) {
			f.Value.Set(f.DefValue)
			f.Changed = false
		})
	})
}

func TestCreateRimWithValidFile(t *testing.T) {
	// Arrange
	setupCreateRimTest(t)
	server := mockServer(t)
	defer server.Close()
	setupMockConfiguration(server.URL, tempConfigFile)

	rimFilePath := "../test/resources/rim-policy.txt"

	// Act - Test 1: Create RIM with valid file
	args := []string{constants.CreateCmd, constants.RimCmd, "-n", "test.rim.name", "-f", rimFilePath}
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err)

	// Act - Test 2: Create RIM with request ID
	args2 := []string{constants.CreateCmd, constants.RimCmd, "-n", "test.rim2", "-f", rimFilePath, "-q", "valid-request-id"}
	_, err = execute(t, tenantCmd, args2)

	// Assert
	assert.NoError(t, err)
}

func TestCreateRimWithNoParams(t *testing.T) {
	// Arrange
	setupCreateRimTest(t)
	args := []string{constants.CreateCmd, constants.RimCmd}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "required flag(s) \"rim-content-file\", \"rim-name\" not set")
}

func TestCreateRimWithNoContentFile(t *testing.T) {
	// Arrange
	setupCreateRimTest(t)
	args := []string{constants.CreateCmd, constants.RimCmd, "-n", "test.rim.name"}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "required flag(s) \"rim-content-file\" not set")
}

func TestCreateRimWithNoName(t *testing.T) {
	// Arrange
	setupCreateRimTest(t)
	args := []string{constants.CreateCmd, constants.RimCmd, "-f", "../test/resources/rim-policy.txt"}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "required flag(s) \"rim-name\" not set")
}

func TestCreateRimWithInvalidFilePath(t *testing.T) {
	// Arrange
	setupCreateRimTest(t)
	args := []string{constants.CreateCmd, constants.RimCmd, "-n", "test.rim", "-f", "/nonexistent/path/file.json"}

	// Act
	_, err := execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "File does not exist")
}

func TestCreateRimWithInvalidJSON(t *testing.T) {
	// Arrange
	setupCreateRimTest(t)
	tmpFile, err := os.CreateTemp("", "invalid-rim-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("{invalid json")
	assert.NoError(t, err)
	tmpFile.Close()

	args := []string{constants.CreateCmd, constants.RimCmd, "-n", "test.rim", "-f", tmpFile.Name()}

	// Act
	_, err = execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "RIM content must be either valid JSON or valid JWT format")
}

func TestCreateRimWithInvalidRequestId(t *testing.T) {
	// Arrange
	setupCreateRimTest(t)
	args := []string{constants.CreateCmd, constants.RimCmd, "-n", "test.rim", "-f", "somefile", "-q", "@#$invalid-id"}

	// Act
	_, err = execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "Request ID should be at most 128 characters long and "+
		"should contain only alphanumeric characters, _, space, - or /")
}
func TestCreateRimWithInvalidName(t *testing.T) {
	// Arrange
	setupCreateRimTest(t)
	tmpFile, err := os.CreateTemp("", "test-rim-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(`{"test": "data"}`)
	assert.NoError(t, err)
	tmpFile.Close()

	testCases := []struct {
		name        string
		rimName     string
		expectedErr string
	}{
		{
			name:        "Name with special characters",
			rimName:     "$$$$$a$gbccc",
			expectedErr: "RIM name is invalid",
		},
		{
			name:        "Name with spaces",
			rimName:     "test rim name",
			expectedErr: "RIM name is invalid",
		},
		{
			name:        "Name starting with special character",
			rimName:     "-testrim",
			expectedErr: "RIM name is invalid",
		},
		{
			name:        "Name ending with special character",
			rimName:     "testrim-",
			expectedErr: "RIM name is invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := []string{constants.CreateCmd, constants.RimCmd, "-n", tc.rimName, "-f", tmpFile.Name()}

			// Act
			_, err := execute(t, tenantCmd, args)

			// Assert
			assert.ErrorContains(t, err, tc.expectedErr)
		})
	}
}
func TestCreateRimWithInvalidUrl(t *testing.T) {
	// Arrange
	setupCreateRimTest(t)
	tmpFile, err := os.CreateTemp("", "valid-rim-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(`{"measurements":[]}`)
	assert.NoError(t, err)
	tmpFile.Close()

	load, err := config.LoadConfiguration()
	assert.NoError(t, err)

	viper.Set("trustauthority-url", "bogus\nbase\nURL")
	defer viper.Set("trustauthority-url", load.TrustAuthorityBaseUrl)

	args := []string{constants.CreateCmd, constants.RimCmd, "-n", "test.rim", "-f", tmpFile.Name()}

	// Act
	_, err = execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "invalid control character in URL")
}

func TestCreateRimWithSignedJWT(t *testing.T) {
	// Arrange
	setupCreateRimTest(t)
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

	args := []string{constants.CreateCmd, constants.RimCmd, "-n", "test.rim.signed", "-f", tmpFile.Name()}

	// Act
	_, err = execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err, "Should accept signed JWT content")
}

func TestCreateRimWithUnsignedJWT(t *testing.T) {
	// Arrange
	setupCreateRimTest(t)
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

	args := []string{constants.CreateCmd, constants.RimCmd, "-n", "test.rim.unsigned", "-f", tmpFile.Name()}

	// Act
	_, err = execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err, "Should accept unsigned JWT content (alg=none)")
}

func TestCreateRimWithJWTAndWhitespace(t *testing.T) {
	// Arrange
	setupCreateRimTest(t)
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

	args := []string{constants.CreateCmd, constants.RimCmd, "-n", "test.rim.whitespace", "-f", tmpFile.Name()}

	// Act
	_, err = execute(t, tenantCmd, args)

	// Assert
	assert.NoError(t, err, "Should handle JWT with whitespace")
}

func TestCreateRimWithInvalidJWT(t *testing.T) {
	// Arrange
	setupCreateRimTest(t)
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

	args := []string{constants.CreateCmd, constants.RimCmd, "-n", "test.rim.invalid", "-f", tmpFile.Name()}

	// Act
	_, err = execute(t, tenantCmd, args)

	// Assert
	assert.ErrorContains(t, err, "RIM content is not valid JSON or JWT format")
}

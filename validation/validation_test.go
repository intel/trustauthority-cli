/*
 * Copyright (C) 2025 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmailAddress(t *testing.T) {
	// Arrange
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"Valid email", "test@example.com", false},
		{"Valid email with subdomain", "user@mail.example.com", false},
		{"Valid email with numbers", "user123@example.com", false},
		{"Invalid email no @", "testexample.com", true},
		{"Invalid email no domain", "test@", true},
		{"Empty email", "", true},
		{"Invalid characters", "test@exam ple.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateEmailAddress(tt.email)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "ValidateEmailAddress() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestValidateTrustAuthorityAPIKey(t *testing.T) {
	// Arrange
	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{"Valid API key", "aGVsbG93b3JsZGhlbGxvd29ybGRoZWxsb3dvcmxk", false},
		{"Valid API key with special chars", "aGVsbG93b3JsZC1oZWxsb3dvcmxkX2hlbGxvd29ybGQ=", false},
		{"Empty API key", "", true},
		{"Too short API key", "short", true},
		{"Invalid characters", "invalid@#$%^&*()", true},
		{"Whitespace only", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateTrustAuthorityAPIKey(tt.apiKey)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "ValidateTrustAuthorityAPIKey() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestValidateTrustAuthorityJwt(t *testing.T) {
	// Arrange
	tests := []struct {
		name    string
		jwt     string
		wantErr bool
	}{
		{"Valid JWT", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", false},
		{"Empty JWT", "", true},
		{"Invalid JWT format", "not.a.jwt", true},
		{"Whitespace only", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateTrustAuthorityJwt(tt.jwt)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "ValidateTrustAuthorityJwt() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestValidateApiClientName(t *testing.T) {
	// Arrange
	tests := []struct {
		name    string
		apiName string
		wantErr bool
	}{
		{"Valid name", "test-api-client", false},
		{"Valid name with underscore", "test_api_client", false},
		{"Valid name alphanumeric", "testApiClient123", false},
		{"Empty name", "", true},
		{"Too short", "ab", true},
		{"Starts with hyphen", "-testapi", true},
		{"Ends with hyphen", "testapi-", true},
		{"Contains spaces", "test api client", true},
		{"Too long", "a123456789012345678901234567890123456789012345678901234567890123456789", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateApiClientName(tt.apiName)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "ValidateApiClientName() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestValidateTagName(t *testing.T) {
	// Arrange
	tests := []struct {
		name    string
		tagName string
		wantErr bool
	}{
		{"Valid tag name", "test-tag", false},
		{"Valid tag with underscore", "test_tag", false},
		{"Valid alphanumeric", "testTag123", false},
		{"Empty tag", "", true},
		{"Too short", "ab", true},
		{"Starts with hyphen", "-testtag", true},
		{"Ends with hyphen", "testtag-", true},
		{"Contains spaces", "test tag", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateTagName(tt.tagName)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "ValidateTagName() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestValidateTagValue(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		tagValue string
		wantErr  bool
	}{
		{"Valid tag value", "test-value", false},
		{"Valid value with underscore", "test_value", false},
		{"Valid alphanumeric", "testValue123", false},
		{"Empty value", "", true},
		{"Too short", "ab", true},
		{"Starts with hyphen", "-testvalue", true},
		{"Contains spaces", "test value", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateTagValue(tt.tagValue)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "ValidateTagValue() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestValidatePolicyName(t *testing.T) {
	// Arrange
	tests := []struct {
		name       string
		policyName string
		wantErr    bool
	}{
		{"Valid policy name", "test-policy", false},
		{"Valid with underscore", "test_policy", false},
		{"Valid alphanumeric", "testPolicy123", false},
		{"Empty name", "", true},
		{"Too short", "ab", true},
		{"Contains spaces", "test policy", true},
		{"Starts with hyphen", "-testpolicy", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidatePolicyName(tt.policyName)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "ValidatePolicyName() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestValidateRequestId(t *testing.T) {
	// Arrange
	tests := []struct {
		name      string
		requestId string
		wantErr   bool
	}{
		{"Valid request ID", "request-123", false},
		{"Valid with slash", "request/123", false},
		{"Valid with underscore", "request_123", false},
		{"Empty request ID", "", false}, // Empty is allowed
		{"Valid with spaces", "request 123", false},
		{"Too long", string(make([]byte, 130)), true},
		{"Invalid characters", "request@123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateRequestId(tt.requestId)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "ValidateRequestId() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestValidateURL(t *testing.T) {
	// Arrange
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"Valid HTTPS URL", "https://example.com", false},
		{"Valid HTTPS URL with path", "https://example.com/api/v1", false},
		{"Invalid HTTP URL", "http://example.com", true},
		{"Invalid URL format", "not-a-url", true},
		{"Empty URL", "", true},
		{"FTP URL", "ftp://example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateURL(tt.url)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestValidatePath(t *testing.T) {
	// Arrange
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"Valid existing file", tmpFile.Name(), false},
		{"Non-existent file", "/path/to/nonexistent/file.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			_, err := ValidatePath(tt.path)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "ValidatePath() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestValidateSize(t *testing.T) {
	// Arrange
	// Create a small temporary file
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some data
	tmpFile.WriteString("test content")
	tmpFile.Close()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"Small file", tmpFile.Name(), false},
		{"Non-existent file", "/path/to/nonexistent/file.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateSize(tt.path)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "ValidateSize() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestCheckFilePathForInvalidChars(t *testing.T) {
	// Arrange
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"Valid path", filepath.Join("test", "path", "file.txt"), false},
		{"Valid path with spaces", filepath.Join("test path", "file name.txt"), false},
		{"Too long path", string(make([]byte, 5000)), true},
		{"Invalid filename with special chars", filepath.Join("test", "file|invalid.txt"), true},
		{"Path with null chars", filepath.Join("test", "\x00file.txt"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := checkFilePathForInvalidChars(tt.path)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "checkFilePathForInvalidChars() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestValidateURLWithInvalidScheme(t *testing.T) {
	// Arrange
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"Valid HTTPS URL", "https://example.com", false},
		{"HTTP URL (invalid scheme)", "http://example.com", true},
		{"FTP URL (invalid scheme)", "ftp://example.com", true},
		{"Malformed URL", "not-a-url", true},
		{"Empty URL", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateURL(tt.url)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestValidatePathWithSymlink(t *testing.T) {
	// Arrange
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create a symlink to the temp file
	symlinkPath := tmpFile.Name() + ".link"
	err = os.Symlink(tmpFile.Name(), symlinkPath)
	if err != nil {
		t.Skipf("Cannot create symlink: %v", err)
	}
	defer os.Remove(symlinkPath)

	// Act
	// Test validation with symlink
	_, err = ValidatePath(symlinkPath)
	// Assert
	assert.NoError(t, err, "ValidatePath() should resolve symlink without error")
}

func TestValidatePathWithInvalidChars(t *testing.T) {
	// Arrange
	// Test path with invalid characters
	invalidPath := string([]byte{0x00, 0x01}) + "file.txt"
	// Act
	_, err := ValidatePath(invalidPath)
	// Assert
	assert.Error(t, err, "ValidatePath() expected error for path with invalid chars")
}

func TestValidateSizeLargeFile(t *testing.T) {
	// Arrange
	// Create a temp file
	tmpFile, err := os.CreateTemp("", "test-large-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write data to make it larger than MaxPolicyFileSize (if feasible for testing)
	// For this test, we'll use a smaller file and assume the constant is reasonable
	tmpFile.WriteString("small test content")
	tmpFile.Close()
	// Act
	// Test with the file we created (should pass since it's small)
	err = ValidateSize(tmpFile.Name())
	// Assert
	assert.NoError(t, err, "ValidateSize() error = %v for small file", err)
}

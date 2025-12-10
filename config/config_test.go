/*
 * Copyright (C) 2025 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration(t *testing.T) {
	// Arrange
	tmpDir, cleanupTmp := setupTestEnv(t)
	defer cleanupTmp()

	_, cleanupHome := setupHomeDir(tmpDir)
	defer cleanupHome()

	configDir := createConfigDir(tmpDir, t)
	configContent := `trustauthority-url: https://test.example.com
trustauthority-api-key: testApiKey12345678901234567890
log-level: info
http-client-timeout: 30
`
	createConfigFile(configDir, configContent, t)

	// Reset viper for this test
	viper.Reset()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(configDir)

	t.Run("Load valid configuration", func(t *testing.T) {
		// Act
		config, err := LoadConfiguration()
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "https://test.example.com", config.TrustAuthorityBaseUrl)
	})

	// Test with non-existent config
	viper.Reset()
	viper.SetConfigName("nonexistent")
	viper.SetConfigType("yml")
	viper.AddConfigPath(tmpDir)

	t.Run("Load non-existent configuration", func(t *testing.T) {
		// Act
		_, err := LoadConfiguration()
		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Config File \"nonexistent\" Not Found")
	})
}

func TestSetupConfigWithValidEnvVars(t *testing.T) {
	// Arrange
	tmpDir, cleanupTmp := setupTestEnv(t)
	defer cleanupTmp()

	// Create env file with only allowed env vars
	envContent := `TRUSTAUTHORITY_URL=https://test.example.com
TRUSTAUTHORITY_API_KEY=dGVzdEFwaUtleTEyMzQ1Njc4OTAxMjM0NTY3ODkw
`
	envFilePath := createEnvFile(tmpDir, envContent, t)

	testHome, cleanupHome := setupHomeDir(tmpDir, "TRUSTAUTHORITY_URL", "TRUSTAUTHORITY_API_KEY")
	defer cleanupHome()

	createConfigDir(testHome, t)

	viper.Reset()
	// Act
	err := SetupConfig(envFilePath)
	// Assert
	assert.NoError(t, err)
}

func TestSetupConfigMissingURL(t *testing.T) {
	// Arrange
	tmpDir, cleanupTmp := setupTestEnv(t)
	defer cleanupTmp()

	// Create env file without TRUSTAUTHORITY_URL
	envContent := `TRUSTAUTHORITY_API_KEY=testApiKey12345678901234567890
`
	envFilePath := createEnvFile(tmpDir, envContent, t)

	// Setup home directory
	testHome, cleanupHome := setupHomeDir(tmpDir)
	defer cleanupHome()

	createConfigDir(testHome, t)

	viper.Reset()
	// Act
	err := SetupConfig(envFilePath)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Trust Authority base URL needs to be provided in configuration")
}

func TestSetupConfigInvalidURL(t *testing.T) {
	// Arrange
	tmpDir, cleanupTmp := setupTestEnv(t)
	defer cleanupTmp()

	// Create env file with invalid URL (http instead of https)
	envContent := `TRUSTAUTHORITY_URL=http://test.example.com
TRUSTAUTHORITY_API_KEY=testApiKey12345678901234567890
`
	envFilePath := createEnvFile(tmpDir, envContent, t)

	// Setup home directory
	testHome, cleanupHome := setupHomeDir(tmpDir)
	defer cleanupHome()

	createConfigDir(testHome, t)

	viper.Reset()
	// Act
	err := SetupConfig(envFilePath)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Invalid Trust Authority base URL, URL scheme must be https")
}

func TestSetupConfigMissingAPIKey(t *testing.T) {
	// Arrange
	tmpDir, cleanupTmp := setupTestEnv(t)
	defer cleanupTmp()

	// Create env file without API key
	envContent := `TRUSTAUTHORITY_URL=https://test.example.com
`
	envFilePath := createEnvFile(tmpDir, envContent, t)

	// Setup home directory and ensure TRUSTAUTHORITY_API_KEY is unset
	testHome, cleanupHome := setupHomeDir(tmpDir, "TRUSTAUTHORITY_API_KEY")
	defer cleanupHome()
	os.Unsetenv("TRUSTAUTHORITY_API_KEY")

	createConfigDir(testHome, t)

	viper.Reset()
	// Act
	err := SetupConfig(envFilePath)
	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Trust Authority API Key needs to be provided in configuration")
}

func TestSetupConfigWithJWT(t *testing.T) {
	// Arrange
	tmpDir, cleanupTmp := setupTestEnv(t)
	defer cleanupTmp()

	// Create env file with valid JWT token
	envContent := `TRUSTAUTHORITY_URL=https://test.example.com
TRUSTAUTHORITY_API_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
`
	envFilePath := createEnvFile(tmpDir, envContent, t)

	testHome, cleanupHome := setupHomeDir(tmpDir, "TRUSTAUTHORITY_URL", "TRUSTAUTHORITY_API_KEY")
	defer cleanupHome()

	createConfigDir(testHome, t)

	viper.Reset()
	// Act
	err := SetupConfig(envFilePath)
	// Assert
	assert.NoError(t, err)
}

func TestSetupConfigWithCustomLogLevel(t *testing.T) {
	// Arrange
	tmpDir, cleanupTmp := setupTestEnv(t)
	defer cleanupTmp()

	envContent := `TRUSTAUTHORITY_URL=https://test.example.com
TRUSTAUTHORITY_API_KEY=testApiKey12345678901234567890
`
	envFilePath := createEnvFile(tmpDir, envContent, t)

	testHome, cleanupHome := setupHomeDir(tmpDir, "TRUSTAUTHORITY_URL", "TRUSTAUTHORITY_API_KEY", "LOG_LEVEL", "HTTP_CLIENT_TIMEOUT")
	defer cleanupHome()

	// Set log level and timeout via environment
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("HTTP_CLIENT_TIMEOUT", "60")

	createConfigDir(testHome, t)

	viper.Reset()
	// Act
	err := SetupConfig(envFilePath)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "debug", viper.GetString("log-level"))
}

func TestSetupConfigInvalidLogLevel(t *testing.T) {
	// Arrange
	tmpDir, cleanupTmp := setupTestEnv(t)
	defer cleanupTmp()

	envContent := `TRUSTAUTHORITY_URL=https://test.example.com
TRUSTAUTHORITY_API_KEY=testApiKey12345678901234567890
`
	envFilePath := createEnvFile(tmpDir, envContent, t)

	testHome, cleanupHome := setupHomeDir(tmpDir, "TRUSTAUTHORITY_URL", "TRUSTAUTHORITY_API_KEY", "LOG_LEVEL")
	defer cleanupHome()

	// Set invalid log level via environment
	os.Setenv("LOG_LEVEL", "invalid-level")

	createConfigDir(testHome, t)

	viper.Reset()
	// Act
	// Should not error - invalid log level should just use default
	err := SetupConfig(envFilePath)
	// Assert
	assert.NoError(t, err)
}

func TestLoadConfigurationUnmarshalError(t *testing.T) {
	// Arrange
	tmpDir, cleanupTmp := setupTestEnv(t)
	defer cleanupTmp()

	_, cleanupHome := setupHomeDir(tmpDir)
	defer cleanupHome()

	configDir := createConfigDir(tmpDir, t)

	// Create an invalid YAML config file (malformed structure)
	configContent := `trustauthority-url: https://test.example.com
trustauthority-api-key: [invalid, array, structure]
http-client-timeout: "not_a_number"
`
	createConfigFile(configDir, configContent, t)

	viper.Reset()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(configDir)

	// Act
	config, err := LoadConfiguration()
	// Assert
	// Should handle gracefully even with type mismatches
	assert.NotNil(t, config)
	assert.Error(t, err, "Failed to unmarshal config")
}

func TestSetupConfigInvalidAPIKey(t *testing.T) {
	// Arrange
	tmpDir, cleanupTmp := setupTestEnv(t)
	defer cleanupTmp()

	envContent := `TRUSTAUTHORITY_URL=https://test.example.com
TRUSTAUTHORITY_API_KEY=invalid-key-too-short
`
	envFilePath := createEnvFile(tmpDir, envContent, t)

	testHome, cleanupHome := setupHomeDir(tmpDir, "TRUSTAUTHORITY_URL", "TRUSTAUTHORITY_API_KEY")
	defer cleanupHome()

	createConfigDir(testHome, t)

	viper.Reset()
	// Act
	err := SetupConfig(envFilePath)
	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Invalid Trust Authority Api key")
}

func TestLoadConfigurationWithInvalidTimeout(t *testing.T) {
	// Arrange
	tmpDir, cleanupTmp := setupTestEnv(t)
	defer cleanupTmp()

	_, cleanupHome := setupHomeDir(tmpDir)
	defer cleanupHome()

	configDir := createConfigDir(tmpDir, t)

	// Create a config file with invalid structure
	configContent := `trustauthority-url: https://test.example.com
trustauthority-api-key: testKey
log-level: info
http-client-timeout: not_a_number
`
	createConfigFile(configDir, configContent, t)

	viper.Reset()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(configDir)

	// Act
	config, err := LoadConfiguration()
	// Assert
	// Should still succeed as viper handles type conversion
	assert.NotNil(t, config)
	assert.Error(t, err, "Failed to unmarshal config")
}

// setupTestEnv creates a temporary test environment and returns cleanup function
func setupTestEnv(t *testing.T) (tmpDir string, cleanup func()) {
	tmpDir, err := os.MkdirTemp("", "config-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	cleanup = func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// ------------ Test Helpers ------------
// setupHomeDir sets up a test home directory and returns cleanup function
func setupHomeDir(tmpDir string, envVarsToClean ...string) (testHome string, cleanup func()) {
	testHome = filepath.Join(tmpDir, "home")
	os.MkdirAll(testHome, 0755)

	originalHome := os.Getenv("HOME")
	if originalHome == "" {
		originalHome = os.Getenv("USERPROFILE")
	}

	os.Setenv("HOME", testHome)
	os.Setenv("USERPROFILE", testHome)

	cleanup = func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
			os.Setenv("USERPROFILE", originalHome)
		}
		for _, envVar := range envVarsToClean {
			os.Unsetenv(envVar)
		}
	}

	return testHome, cleanup
}

// createEnvFile creates an environment file with the given content
func createEnvFile(tmpDir, content string, t *testing.T) string {
	envFilePath := filepath.Join(tmpDir, "test.env")
	err := os.WriteFile(envFilePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write env file: %v", err)
	}
	return envFilePath
}

// createConfigDir creates the config directory and returns its path
func createConfigDir(testHome string, t *testing.T) string {
	configDir := filepath.Join(testHome, ".config", "trustauthorityctl")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}
	return configDir
}

// createConfigFile creates a config file with the given content
func createConfigFile(configDir, content string, t *testing.T) string {
	configPath := filepath.Join(configDir, "config.yml")
	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
	return configPath
}

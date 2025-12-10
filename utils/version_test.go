/*
 * Copyright (C) 2025 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package utils

import (
	"encoding/json"
	"testing"
)

func TestGetVersion(t *testing.T) {
	version, err := GetVersion()
	if err != nil {
		t.Errorf("GetVersion() error = %v", err)
	}

	if version == "" {
		t.Error("GetVersion() returned empty string")
	}

	// Verify it's valid JSON
	var v CliVersion
	err = json.Unmarshal([]byte(version), &v)
	if err != nil {
		t.Errorf("GetVersion() returned invalid JSON: %v", err)
	}

	// Verify structure
	if v.Name == "" {
		t.Error("GetVersion() Name field is empty")
	}
}

func TestCliVersion(t *testing.T) {
	// Test that the cliVersion variable is properly initialized
	if cliVersion.Name == "" {
		t.Error("cliVersion.Name is empty")
	}
}

func TestVersionVariables(t *testing.T) {
	// These variables are set at build time, so they may be empty in tests
	// We just verify they're accessible
	_ = Version
	_ = GitHash
	_ = BuildDate

	// Test the structure can be marshaled
	testVersion := CliVersion{
		Name:      "Test CLI",
		Version:   "1.0.0",
		GitHash:   "abc123",
		BuildDate: "2025-01-01",
	}

	data, err := json.Marshal(testVersion)
	if err != nil {
		t.Errorf("Failed to marshal CliVersion: %v", err)
	}

	if len(data) == 0 {
		t.Error("Marshaled CliVersion is empty")
	}

	// Verify it contains expected fields
	if !contains(string(data), "name") {
		t.Error("Marshaled data missing 'name' field")
	}
}

func contains(str, substr string) bool {
	return len(str) > 0 && len(substr) > 0 && str != substr &&
		(str[:len(substr)] == substr ||
			(len(str) > len(substr) && contains(str[1:], substr)))
}

func TestGetVersionFields(t *testing.T) {
	// Test that GetVersion returns all expected JSON fields
	version, err := GetVersion()
	if err != nil {
		t.Errorf("GetVersion() error = %v", err)
	}

	// Parse the JSON to verify structure
	var v map[string]interface{}
	err = json.Unmarshal([]byte(version), &v)
	if err != nil {
		t.Errorf("GetVersion() returned invalid JSON: %v", err)
	}

	// Verify all fields are present
	expectedFields := []string{"name", "version", "gitHash", "buildDate"}
	for _, field := range expectedFields {
		if _, ok := v[field]; !ok {
			t.Errorf("GetVersion() missing field: %s", field)
		}
	}
}

func TestGetVersionWithAllFields(t *testing.T) {
	// Save original values
	origVersion := Version
	origGitHash := GitHash
	origBuildDate := BuildDate
	origCliVersion := cliVersion

	defer func() {
		// Restore original values
		Version = origVersion
		GitHash = origGitHash
		BuildDate = origBuildDate
		cliVersion = origCliVersion
	}()

	// Set test values
	Version = "1.2.3"
	GitHash = "abc123def456"
	BuildDate = "2025-12-08"
	cliVersion = CliVersion{
		Name:      "Test CLI",
		Version:   Version,
		GitHash:   GitHash,
		BuildDate: BuildDate,
	}

	version, err := GetVersion()
	if err != nil {
		t.Fatalf("GetVersion() error = %v", err)
	}

	// Verify all fields in JSON
	var v CliVersion
	err = json.Unmarshal([]byte(version), &v)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if v.Name != "Test CLI" {
		t.Errorf("Name = %v, want Test CLI", v.Name)
	}
	if v.Version != "1.2.3" {
		t.Errorf("Version = %v, want 1.2.3", v.Version)
	}
	if v.GitHash != "abc123def456" {
		t.Errorf("GitHash = %v, want abc123def456", v.GitHash)
	}
	if v.BuildDate != "2025-12-08" {
		t.Errorf("BuildDate = %v, want 2025-12-08", v.BuildDate)
	}
}

func TestGetVersionEmptyFields(t *testing.T) {
	// Test with empty fields (like in development builds)
	origCliVersion := cliVersion

	defer func() {
		cliVersion = origCliVersion
	}()

	cliVersion = CliVersion{
		Name:      "Test",
		Version:   "",
		GitHash:   "",
		BuildDate: "",
	}

	version, err := GetVersion()
	if err != nil {
		t.Fatalf("GetVersion() should not error with empty fields: %v", err)
	}

	if version == "" {
		t.Error("GetVersion() should not return empty string")
	}

	// Verify it's valid JSON even with empty fields
	var v CliVersion
	err = json.Unmarshal([]byte(version), &v)
	if err != nil {
		t.Errorf("Invalid JSON with empty fields: %v", err)
	}
}

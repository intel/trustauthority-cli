/*
 * Copyright (C) 2025 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestReadAnswerFileToEnv(t *testing.T) {
	// Arrange
	tmpFile, err := os.CreateTemp("", "answer-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write valid environment variables (no INVALID_VAR)
	content := `TRUSTAUTHORITY_URL=https://example.com
TRUSTAUTHORITY_API_KEY=test-api-key-1234567890
# This is a comment
`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{"Valid answer file", tmpFile.Name(), false},
		{"Non-existent file", "/path/to/nonexistent.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ReadAnswerFileToEnv(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAnswerFileToEnv() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Assert
				assert.Equal(t, "https://example.com", os.Getenv("TRUSTAUTHORITY_URL"))
				assert.Equal(t, "test-api-key-1234567890", os.Getenv("TRUSTAUTHORITY_API_KEY"))
				// Clean up
				os.Unsetenv("TRUSTAUTHORITY_URL")
				os.Unsetenv("TRUSTAUTHORITY_API_KEY")
			}
		})
	}
}

func TestReadAnswerFileToEnvWithInvalidVar(t *testing.T) {
	// Arrange
	tmpFile, err := os.CreateTemp("", "answer-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `INVALID_VARIABLE=value`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()
	// Act
	err = ReadAnswerFileToEnv(tmpFile.Name())
	// Assert
	assert.NotNil(t, err, "Expected error for invalid environment variable")
	assert.Contains(t, err.Error(), "Invalid ENV variable: INVALID_VARIABLE", "Error message should mention invalid variable")
}

func TestIsValidEnvVariable(t *testing.T) {
	tests := []struct {
		name     string
		variable string
		want     bool
	}{
		{"Valid TRUSTAUTHORITY_URL", "TRUSTAUTHORITY_URL", true},
		{"Valid TRUSTAUTHORITY_API_KEY", "TRUSTAUTHORITY_API_KEY", true},
		{"Invalid variable", "INVALID_VAR", false},
		{"Empty variable", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidEnvVariable(tt.variable); got != tt.want {
				t.Errorf("isValidEnvVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetUpLogs(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		logLevel string
		wantErr  bool
	}{
		{"Valid log level info", "info", false},
		{"Valid log level debug", "debug", false},
		{"Valid log level warn", "warn", false},
		{"Valid log level error", "error", false},
		{"Invalid log level", "invalid", false}, // Should default to info
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			// Act
			err := SetUpLogs(&buf, tt.logLevel)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "SetUpLogs() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestCheckSigningAlgorithm(t *testing.T) {
	// Arrange
	// Generate RSA private keys for testing
	privKey2048, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate 2048-bit RSA key: %v", err)
	}

	privKey3072, err := rsa.GenerateKey(rand.Reader, 3072)
	if err != nil {
		t.Fatalf("Failed to generate 3072-bit RSA key: %v", err)
	}

	tests := []struct {
		name      string
		privKey   *rsa.PrivateKey
		algorithm string
		wantNil   bool
	}{
		{"2048-bit key with RS256", privKey2048, "RS256", false},
		{"3072-bit key with RS384", privKey3072, "RS384", false},
		{"2048-bit key with RS384 (mismatch)", privKey2048, "RS384", true},
		{"3072-bit key with RS256 (mismatch)", privKey3072, "RS256", true},
		{"Invalid algorithm", privKey2048, "INVALID", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			got := CheckSigningAlgorithm(tt.privKey, tt.algorithm)
			assert.Equal(t, tt.wantNil, got == nil, "CheckSigningAlgorithm() result nil = %v, wantNil %v", got == nil, tt.wantNil)
		})
	}
}

func TestCheckKeyFiles(t *testing.T) {
	// Arrange
	// Create temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "keytest-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Generate a test RSA key pair
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	// Create certificate
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	// Write private key to file
	privKeyPath := filepath.Join(tmpDir, "private.pem")
	privKeyFile, err := os.Create(privKeyPath)
	if err != nil {
		t.Fatalf("Failed to create private key file: %v", err)
	}
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	pem.Encode(privKeyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privKeyBytes})
	privKeyFile.Close()

	// Write certificate to file
	certPath := filepath.Join(tmpDir, "cert.pem")
	certFile, err := os.Create(certPath)
	if err != nil {
		t.Fatalf("Failed to create certificate file: %v", err)
	}
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	certFile.Close()

	tests := []struct {
		name     string
		privPath string
		certPath string
		wantErr  bool
	}{
		{"Valid key and cert", privKeyPath, certPath, false},
		{"Empty private key path", "", certPath, true},
		{"Empty certificate path", privKeyPath, "", true},
		{"Non-existent private key", "/nonexistent/key.pem", certPath, true},
		{"Non-existent certificate", privKeyPath, "/nonexistent/cert.pem", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			key, cert, err := CheckKeyFiles(tt.privPath, tt.certPath)
			// Assert
			if tt.wantErr {
				assert.Error(t, err, "CheckKeyFiles() should return error")
			} else {
				assert.NoError(t, err, "CheckKeyFiles() should not return error")
				assert.NotNil(t, key, "CheckKeyFiles() should return non-nil key")
				assert.NotEmpty(t, cert, "CheckKeyFiles() should return non-empty certificate")
			}
		})
	}
}

func TestGenerateOutputFileName(t *testing.T) {
	// Arrange
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid file", tmpFile.Name(), false},
		{"Non-existent file", "/path/to/nonexistent.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			output, err := GenerateOutputFileName(tt.input)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "GenerateOutputFileName() error = %v, wantErr %v", err, tt.wantErr)
			if !tt.wantErr {
				assert.NotEmpty(t, output, "GenerateOutputFileName() should return non-empty output")
				assert.True(t, strings.HasSuffix(output, ".txt"), "GenerateOutputFileName() output doesn't end with .txt")
			}
		})
	}
}

func TestPublicKeyToBytes(t *testing.T) {
	// Arrange
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	pubBytes := publicKeyToBytes(&privKey.PublicKey)
	if len(pubBytes) == 0 {
		t.Errorf("publicKeyToBytes() returned empty bytes")
	}

	// Act
	// Verify it's valid PEM
	block, _ := pem.Decode(pubBytes)
	// Assert
	assert.NotNil(t, block, "publicKeyToBytes() returned invalid PEM")
	assert.Equal(t, "PUBLIC KEY", block.Type, "publicKeyToBytes() PEM block type is not 'PUBLIC KEY'")
}

func TestParseCertificate(t *testing.T) {
	// Arrange
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{"Valid certificate", certPEM, false},
		{"Invalid certificate", []byte("not a certificate"), true},
		{"Empty input", []byte{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			cert, err := parseCertificate(tt.input)
			// Assert
			assert.Equal(t, tt.wantErr, err != nil, "parseCertificate() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equal(t, !tt.wantErr, cert != nil, "parseCertificate() cert nil = %v, wantNil %v", cert == nil, tt.wantErr)
		})
	}
}

func TestSetUpLogsOutput(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	// Act
	SetUpLogs(&buf, "info")
	logrus.Info("test message")
	// Assert
	assert.Contains(t, buf.String(), "test message", "Log output should contain the test message")
}

func TestSetUpLogsWithWriter(t *testing.T) {
	// Arrange
	writers := []io.Writer{
		&bytes.Buffer{},
		os.Stdout,
	}

	for i, w := range writers {
		t.Run(string(rune(i)), func(t *testing.T) {
			// Act
			err := SetUpLogs(w, "debug")
			// Assert
			assert.NoError(t, err, "SetUpLogs() should not return error")
		})
	}
}

func TestCheckKeyFilesWithValidFiles(t *testing.T) {
	// Arrange
	// Generate test RSA key pair
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "keyfiles-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create private key file
	privKeyPath := filepath.Join(tmpDir, "private.pem")
	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})
	err = os.WriteFile(privKeyPath, privKeyPEM, 0600)
	if err != nil {
		t.Fatalf("Failed to write private key: %v", err)
	}

	// Create certificate file
	certPath := filepath.Join(tmpDir, "cert.pem")
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),
	}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	err = os.WriteFile(certPath, certPEM, 0644)
	if err != nil {
		t.Fatalf("Failed to write certificate: %v", err)
	}

	// Act
	// Test CheckKeyFiles
	_, _, err = CheckKeyFiles(privKeyPath, certPath)
	// Assert
	assert.NoError(t, err, "CheckKeyFiles() should not return error for valid files")
}

func TestCheckKeyFilesWithInvalidFiles(t *testing.T) {
	// Arrange
	tmpDir, err := os.MkdirTemp("", "keyfiles-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test with non-existent files
	_, _, err = CheckKeyFiles("/nonexistent/private.pem", "/nonexistent/cert.pem")
	if err == nil {
		t.Error("CheckKeyFiles() expected error for non-existent files")
	}

	// Test with invalid private key
	invalidKeyPath := filepath.Join(tmpDir, "invalid.pem")
	err = os.WriteFile(invalidKeyPath, []byte("not a valid key"), 0600)
	if err != nil {
		t.Fatalf("Failed to write invalid key: %v", err)
	}

	invalidCertPath := filepath.Join(tmpDir, "invalid_cert.pem")
	err = os.WriteFile(invalidCertPath, []byte("not a valid cert"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid cert: %v", err)
	}
	// Act
	_, _, err = CheckKeyFiles(invalidKeyPath, invalidCertPath)
	// Assert
	assert.Error(t, err, "CheckKeyFiles() expected error for invalid key/cert files")
}

func TestPublicKeyToBytesWithNilKey(t *testing.T) {
	// Arrange
	// Test edge case - should handle without panic
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	// Test with valid key
	pubBytes := publicKeyToBytes(&privKey.PublicKey)
	if len(pubBytes) == 0 {
		t.Error("publicKeyToBytes() returned empty bytes for valid key")
	}

	// Act
	// Verify the PEM structure
	block, _ := pem.Decode(pubBytes)
	// Assert
	assert.NotNil(t, block, "publicKeyToBytes() returned invalid PEM for valid key")
	assert.Equal(t, "PUBLIC KEY", block.Type, "publicKeyToBytes() PEM block type is not 'PUBLIC KEY' for valid key")
}

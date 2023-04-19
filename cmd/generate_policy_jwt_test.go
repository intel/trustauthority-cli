/*
 * Copyright (C) 2023 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/test"
	"math/big"
	"os"
	"testing"
	"time"
)

const (
	keyFile  = "../test/resources/sample-policy-signing-key.txt"
	certFile = "../test/resources/sample-policy-signing-cert.txt"
)

func TestGeneratePolicyJwtCmd(t *testing.T) {
	server := test.MockServer(t)
	defer server.Close()
	test.SetupMockConfiguration(server.URL, tempConfigFile)

	generateKeyPairForTests(t, keyFile, certFile)
	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:    []string{constants.CreateCmd, constants.PolicyJwtCmd},
			wantErr: true,
		},
		{
			args:    []string{constants.CreateCmd, constants.PolicyJwtCmd, "-f", "../test/resources/rego-policy.txt"},
			wantErr: false,
		},
		{
			args:    []string{constants.CreateCmd, constants.PolicyJwtCmd, "-f", "../test/resources/rego-policy.txt", "-s"},
			wantErr: true,
		},
		{
			args: []string{constants.CreateCmd, constants.PolicyJwtCmd, "-f", "../test/resources/rego-policy.txt",
				"-p", keyFile, "-s"},
			wantErr: true,
		},
		{
			args: []string{constants.CreateCmd, constants.PolicyJwtCmd, "-f", "../test/resources/rego-policy.txt",
				"-p", keyFile,
				"-c", certFile, "-s"},
			wantErr: false,
		},
		{
			args: []string{constants.CreateCmd, constants.PolicyJwtCmd, "-f", "../test/resources/rego-policy.txt",
				"-p", keyFile,
				"-c", "..test/", "-s"},
			wantErr:     true,
			description: "Test Invalid certificate file path",
		},
		{
			args: []string{constants.CreateCmd, constants.PolicyJwtCmd, "-f", "../test/resources/rego-policy.txt",
				"-p", "../test",
				"-c", certFile, "-s"},
			wantErr:     true,
			description: "Test Invalid private key path",
		},
		{
			args: []string{constants.CreateCmd, constants.PolicyJwtCmd, "-f", "../test/resources/rego-policy.txt",
				"-p", keyFile,
				"-c", certFile, "-a", constants.RS256},
			wantErr:     true,
			description: "Test Signing algorithm provided as input is not compatible with the private key type",
		},
		{
			args:        []string{constants.CreateCmd, constants.PolicyJwtCmd, "-f", "../test/resources/@rego-policy.txt"},
			wantErr:     true,
			description: "Test Invalid policy file path provided",
		},
		{
			args:        []string{constants.CreateCmd, constants.PolicyJwtCmd, "-f", ""},
			wantErr:     true,
			description: "Test Policy file path cannot be empty",
		},
		{
			args:        []string{constants.CreateCmd, constants.PolicyJwtCmd, "-f", "../test/resources/"},
			wantErr:     true,
			description: "Test Error reading policy file",
		},

		{
			args:        []string{constants.CreateCmd, constants.PolicyJwtCmd, "-f", "../test/resources/rego-policy1.txt"},
			wantErr:     true,
			description: "Test Policy file does not contain a rego policy",
		},
		{
			args: []string{constants.CreateCmd, constants.PolicyJwtCmd, "-f", "../test/resources/rego-policy.txt", "-a",
				"invalid algorithm"},
			wantErr:     true,
			description: "Test Input algorithm is not supported",
		},
	}

	createCmd.AddCommand(createPolicyJwtCmd)
	tenantCmd.AddCommand(createCmd)

	for _, tc := range tt {
		_, err := execute(t, tenantCmd, tc.args)

		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		time.Sleep(1 * time.Second)
	}

	err = os.Remove(keyFile)
	assert.NoError(t, err)

	err = os.Remove(certFile)
	assert.NoError(t, err)
}

func generateKeyPairForTests(t *testing.T, keyFile, certFile string) {
	keyPair, err := rsa.GenerateKey(rand.Reader, 3072)
	assert.NoError(t, err)
	privKeyBytes := x509.MarshalPKCS1PrivateKey(keyPair)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	}
	privatePem, err := os.Create(keyFile)
	assert.NoError(t, err)

	err = pem.Encode(privatePem, privateKeyBlock)
	assert.NoError(t, err)

	keyUsage := x509.KeyUsageDigitalSignature

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Second * 30)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	assert.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Test Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &keyPair.PublicKey, keyPair)
	assert.NoError(t, err)

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}

	publicPem, err := os.Create(certFile)
	assert.NoError(t, err)

	err = pem.Encode(publicPem, certBlock)
	assert.NoError(t, err)
	err = publicPem.Close()
	assert.NoError(t, err)
	err = privatePem.Close()
	assert.NoError(t, err)
}

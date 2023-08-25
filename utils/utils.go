/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package utils

import (
	"bufio"
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"intel/tac/v1/constants"
	models2 "intel/tac/v1/internal/models"
	"intel/tac/v1/validation"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ReadAnswerFileToEnv(filename string) error {
	path, err := validation.ValidatePath(filename)
	if err != nil {
		return errors.Wrap(err, "Invalid ReadAnswerFilePath")
	}
	fin, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "Failed to load answer file")
	}
	scanner := bufio.NewScanner(fin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" ||
			strings.HasPrefix(line, "#") {
			continue
		}
		equalSign := strings.Index(line, "=")
		if equalSign > 0 {
			key := line[0:equalSign]
			val := line[equalSign+1:]
			if key != "" &&
				val != "" {
				err = os.Setenv(key, val)
				if err != nil {
					return errors.Wrap(err, "Failed to set ENV")
				}
			}
			isValid := isValidEnvVariable(key)
			if !isValid {
				return errors.Errorf("Invalid ENV variable: %s", key)
			}
		}
	}
	return nil
}
func isValidEnvVariable(lookup string) bool {
	envMap := map[string]bool{
		"TRUSTAUTHORITY_URL":     true,
		"TRUSTAUTHORITY_API_KEY": true,
	}
	if _, ok := envMap[lookup]; ok {
		return true
	}
	return false
}

// SetUpLogs set the log output and the log level
func SetUpLogs(logFile io.Writer, logLevel string) error {
	logrus.SetOutput(logFile)
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		lvl, _ = logrus.ParseLevel(constants.DefaultLogLevel)
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(lvl)
	}
	return nil
}

// CheckSigningAlgorithm check if provided algorithm makes sense
func CheckSigningAlgorithm(privKeyFinal *rsa.PrivateKey, algorithm string) jwt.SigningMethod {
	if privKeyFinal.N.BitLen() == 2048 && !strings.Contains(algorithm, constants.HashSize256) {
		fmt.Println("Input private key file and algorithm do not match")
		return nil
	}
	if privKeyFinal.N.BitLen() == 3072 && !strings.Contains(algorithm, constants.HashSize384) {
		fmt.Println("Input private key file and algorithm do not match")
		return nil
	}
	signMethod := jwt.GetSigningMethod(algorithm)
	if signMethod == nil {
		fmt.Println("Input signing algorithm not found")
	}
	return signMethod
}

// CheckKeyFiles check input private key and certificate files are valid
func CheckKeyFiles(privKeyFilePath, certificateFilePath string) (*rsa.PrivateKey, string, error) {
	if privKeyFilePath == "" {
		return nil, "", errors.New("Private key file path cannot be empty")
	}

	if certificateFilePath == "" {
		return nil, "", errors.New("Certificate file path cannot be empty")
	}
	filepath, err := validation.ValidatePath(privKeyFilePath)
	if err != nil {
		return nil, "", errors.Wrap(err, "Invalid privKeyFilePath")
	}
	certfile, err := validation.ValidatePath(certificateFilePath)
	if err != nil {
		return nil, "", errors.Wrap(err, "Invalid certificateFilePath")
	}
	privKeyBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, "", errors.Wrap(err, "Error reading private key file")
	}

	privKeyFinal, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes)
	if err != nil {
		return nil, "", errors.Wrap(err, "Error parsing private key PEM file")
	}

	certBytes, err := os.ReadFile(certfile)
	if err != nil {
		return nil, "", errors.Wrap(err, "Error reading certificate file")
	}

	cert, err := parseCertificate(certBytes)
	if err != nil {
		return nil, "", errors.Wrap(err, "Error parsing certificate")
	}

	pubKeyBytesFromCert := publicKeyToBytes(cert.PublicKey.(*rsa.PublicKey))
	pubKeyBytesFromPriv := publicKeyToBytes(&privKeyFinal.PublicKey)

	if bytes.Compare(pubKeyBytesFromCert, pubKeyBytesFromPriv) != 0 {
		return nil, "", errors.New("Provided private key and certificate do not match")
	}

	certContents := base64.StdEncoding.EncodeToString(cert.Raw)
	return privKeyFinal, certContents, nil
}

func GenerateOutputFileName(inputFile string) (string, error) {
	inputFilepath, err := validation.ValidatePath(inputFile)
	if err != nil {
		return "", errors.Wrap(err, "Invalid GenerateOutputFilePath ")
	}
	filename := strings.TrimSuffix(inputFilepath, filepath.Ext(inputFilepath))
	date := time.Now().Format(constants.TimeLayout)

	return filename + ".signed." + date + ".txt", nil
}

func PrintRequestAndTraceId() {
	if models2.RespHeaderFields.RequestId != "" {
		fmt.Println(constants.HTTPHeaderKeyRequestId+": ", models2.RespHeaderFields.RequestId)
	}
	if models2.RespHeaderFields.TraceId != "" {
		fmt.Println(constants.HTTPHeaderKeyTraceId+": ", models2.RespHeaderFields.TraceId)
	}
}

// publicKeyToBytes public key to bytes
func publicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		fmt.Println(err.Error())
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  constants.PublicKey,
		Bytes: pubASN1,
	})

	return pubBytes
}

// parseCertificate parse certificate from unencrypted string format
func parseCertificate(certBytes []byte) (*x509.Certificate, error) {
	certBlock, _ := pem.Decode(certBytes)
	if certBlock == nil || certBlock.Type != constants.CertType {
		return nil, errors.New("not a valid certificate pem file")
	}
	return x509.ParseCertificate(certBlock.Bytes)
}

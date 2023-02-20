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
	"intel/amber/tac/v1/constants"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ReadAnswerFileToEnv(filename string) error {
	fin, err := os.Open(filepath.Clean(filename))
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
		}
	}
	return nil
}

//SetUpLogs set the log output and the log level
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
	privKeyBytes, err := ioutil.ReadFile(filepath.Clean(privKeyFilePath))
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}

	privKeyFinal, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}
	certBytes, err := ioutil.ReadFile(filepath.Clean(certificateFilePath))
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}

	cert, err := parseCertificate(certBytes)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}

	pubKeyBytesFromCert := publicKeyToBytes(cert.PublicKey.(*rsa.PublicKey))
	pubKeyBytesFromPriv := publicKeyToBytes(&privKeyFinal.PublicKey)

	if bytes.Compare(pubKeyBytesFromCert, pubKeyBytesFromPriv) != 0 {
		fmt.Println("Provided private key and certificate do not match")
		return nil, "", errors.New("provided private key and certificate do not match")
	}

	certContents := base64.StdEncoding.EncodeToString(cert.Raw)
	return privKeyFinal, certContents, nil
}

func GenerateOutputFileName(inputFile string) string {
	filename := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	date := time.Now().Format(constants.TimeLayout)
	return filepath.Clean(filename + ".signed." + date + ".txt")
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

/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package validation

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"intel/tac/v1/constants"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	stringReg = regexp.MustCompile("(^[a-zA-Z0-9_ \\/.-]*$)")
	emailReg  = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+\/=?^_'{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)
	// Regex to validate TA API key. Key should contain characters between a-z, A-Z, 0-9, +, /, =, _, -
	// and should be of size between 30 and 128
	apiKeyRegex           = regexp.MustCompile(`^[A-Za-z0-9+/=_-]{30,250}$`)
	subscriptionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-\_]{1,62}[a-zA-Z0-9]$`)
	tagReg                = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-\_]{1,62}[a-zA-Z0-9]$`)
	tagValueReg           = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-\_]{1,62}[a-zA-Z0-9]$`)
	policyNameRegex       = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{1,62}[a-zA-Z0-9]$`)
	requestIdRegex        = regexp.MustCompile(`^[a-zA-Z0-9_ \/.-]{1,128}$`)
	//max length of file name to be allowed in 255 bytes and characters allowed are a-z, A-Z, 0-9, _, ., -
	fileNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_. -]{1,255}$`)
	//in file path, characters allowed are a-z, A-Z, 0-9, _, ., -, \, /, :
	filePathRegex = regexp.MustCompile(`^[a-zA-Z0-9_. :/\\-]*$`)
)

// ValidateStrings method is used to validate input strings
func ValidateStrings(strings []string) error {
	for _, stringValue := range strings {
		if !stringReg.MatchString(stringValue) {
			return errors.New("Invalid string formatted input")
		}
	}
	return nil
}

func ValidateEmailAddress(email string) error {
	if !emailReg.Match([]byte(email)) {
		logrus.Error("Invalid email id provided")
		return errors.New("Invalid email id provided")
	}

	return nil
}

func ValidatePath(path string) (string, error) {
	cleanedPath := filepath.Clean(path)
	if err := checkFilePathForInvalidChars(cleanedPath); err != nil {
		return "", err
	}
	r, err := filepath.EvalSymlinks(cleanedPath)
	if err != nil {
		return "", fmt.Errorf("Unsafe symlink detected in path")
	}
	if err = checkFilePathForInvalidChars(r); err != nil {
		return "", err
	}
	return r, nil
}

func ValidateSize(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.Size() > constants.MaxPolicyFileSize {
		return fmt.Errorf("%s: %d", constants.ErrorInvalidSize, fi.Size())
	}
	return nil
}

func ValidateTrustAuthorityAPIKey(apiKey string) error {
	if strings.TrimSpace(apiKey) == "" {
		return errors.Errorf("%s config variable needs to be set with a proper API Key before using CLI", constants.TrustAuthApiKeyEnvVar)
	}
	if !apiKeyRegex.MatchString(apiKey) {
		return errors.New("Invalid API key found in configuration file. Please update it with a valid API key.")
	}
	return nil
}

func ValidateApiClientName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("ApiClient name cannot be empty")
	}
	if !subscriptionNameRegex.Match([]byte(name)) {
		return errors.New("ApiClient name should be alphanumeric and start with an alphanumeric character with " +
			"_ or - as separator and should be at most 64 characters long")
	}
	return nil
}

func ValidateTagName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("Tag name cannot be empty")
	}
	if !tagReg.Match([]byte(name)) {
		return errors.New("Tag name should be alphanumeric and start with an alphanumeric character with " +
			"_ or - as separator and should be at most 64 characters long")
	}
	return nil
}

func ValidateTagValue(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("Tag value cannot be empty")
	}
	if !tagValueReg.Match([]byte(value)) {
		return errors.New("Tag value should be alphanumeric and start with an alphanumeric character with " +
			"_ or - as separator and should be at most 64 characters long")
	}
	return nil
}

func ValidatePolicyName(policyName string) error {
	if strings.TrimSpace(policyName) == "" {
		return errors.New("Policy name cannot be empty")
	}
	if !policyNameRegex.Match([]byte(policyName)) {
		return errors.New("Policy name is invalid. Policy name should be alpha numeric and have minimum 3 characters with no spaces between words (" +
			"use \"_\" or \"-\" as separators) and should not be more than 64 characters")
	}
	return nil
}

func ValidateRequestId(requestId string) error {
	if strings.TrimSpace(requestId) != "" && !requestIdRegex.Match([]byte(requestId)) {
		return errors.New("Request ID should be at most 128 characters long and should contain only " +
			"alphanumeric characters, _, space, - or \\")
	}
	return nil
}

func ValidateURL(baseURL string) error {
	baseUrl, err := url.Parse(baseURL)
	if err != nil {
		return errors.Wrap(err, "Invalid Trust Authority Base URL")
	}
	if baseUrl.Scheme != constants.HTTPScheme {
		return errors.New("Invalid Trust Authority base URL, URL scheme must be https")
	}
	return nil
}

func checkFilePathForInvalidChars(path string) error {
	filePath, fileName := filepath.Split(path)
	//Max file path length allowed in linux is 4096 characters
	if len(path) > constants.LinuxFilePathSize || !filePathRegex.MatchString(filePath) {
		return errors.New("Invalid linux file path provided")
	}
	if !fileNameRegex.MatchString(fileName) {
		return errors.New("Invalid file name provided")
	}
	return nil
}

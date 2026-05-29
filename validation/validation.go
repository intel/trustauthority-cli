/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package validation

import (
	"encoding/json"
	"fmt"
	"intel/tac/v1/constants"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	emailReg = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+\/=?^_'{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)
	// Regex to validate TA API key. Key should contain characters between a-z, A-Z, 0-9, +, /, =, _, -
	// and should be of size between 30 and 128
	apiKeyRegex           = regexp.MustCompile(`^[A-Za-z0-9+/=_-]{30,250}$`)
	subscriptionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-\_]{1,62}[a-zA-Z0-9]$`)
	tagReg                = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-\_]{1,62}[a-zA-Z0-9]$`)
	tagValueReg           = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-\_]{1,62}[a-zA-Z0-9]$`)
	policyNameRegex       = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{1,62}[a-zA-Z0-9]$`)
	rimNameRegex          = regexp.MustCompile(`^[a-zA-Z](?:[a-zA-Z0-9_.]{1,126}[a-zA-Z0-9])?$`)
	// segmentIdentifierRegex validates each segment as a valid identifier
	// Must start with letter (a-z, A-Z), followed by letters, numbers, or underscores
	// Cannot start with underscore or number, hyphens not allowed
	rimSegmentIdentifierRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	requestIdRegex            = regexp.MustCompile(`^[a-zA-Z0-9_ \/.-]{1,128}$`)
	//max length of file name to be allowed in 255 bytes and characters allowed are a-z, A-Z, 0-9, _, ., -
	fileNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_. -]{1,255}$`)
	//in file path, characters allowed are a-z, A-Z, 0-9, _, ., -, \, /, :
	filePathRegex = regexp.MustCompile(`^[a-zA-Z0-9_. :/\\-]*$`)
)

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

	// Check if file exists first
	if _, err := os.Stat(cleanedPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("File does not exist: %s", cleanedPath)
		}
		return "", fmt.Errorf("Error accessing file: %v", err)
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

func ValidateTrustAuthorityJwt(tokenString string) error {
	if strings.TrimSpace(tokenString) == "" {
		return errors.Errorf("%s config variable needs to be set with a proper API Key before using CLI", constants.TrustAuthApiKeyEnvVar)
	}
	// Parse the token
	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return errors.Wrap(err, "Invalid JWT format found in configuration file. Please update it with a valid JWT.")
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

func ValidateRimName(rimName string) error {
	if strings.TrimSpace(rimName) == "" {
		return errors.New("RIM name cannot be empty")
	}
	if len(rimName) < constants.MinRimNameLength {
		return errors.Errorf("RIM name must be at least %d characters", constants.MinRimNameLength)
	}
	if len(rimName) > constants.MaxRimNameLength {
		return errors.Errorf("RIM name exceeds maximum length of %d characters", constants.MaxRimNameLength)
	}

	// Basic format check: overall pattern validation
	if !rimNameRegex.MatchString(rimName) {
		return errors.New("RIM name is invalid. Name must be 3-128 characters, start with a letter (a-z, A-Z), " +
			"end with alphanumeric (a-z, A-Z, 0-9), and may contain letters, numbers, dots, and underscores. " +
			"Use dots for namespaced names (e.g., acme.rims.mrtd or public.acme.rims.certificates). " +
			"No spaces, hyphens, or special characters allowed. Cannot end with dot or underscore.")
	}

	// Split by dots and validate namespace segments
	segments := strings.Split(rimName, ".")
	if len(segments) > constants.MaxNamespaceSegments {
		return errors.Errorf("RIM name has too many namespace segments. Maximum allowed is %d segments, found %d",
			constants.MaxNamespaceSegments, len(segments))
	}

	// Validate each segment as a valid identifier
	// Each segment must start with a letter (not underscore or number)
	// to ensure it can be used as a code identifier
	for i, segment := range segments {
		if !rimSegmentIdentifierRegex.MatchString(segment) {
			return errors.Errorf("RIM name segment '%s' at position %d is invalid. Each segment must be a valid identifier: "+
				"start with a letter (a-z, A-Z), followed by letters, numbers, or underscores. "+
				"Segments cannot start with a number or underscore, and hyphens are not allowed. "+
				"(e.g., 'public.acme.1.rims' is invalid, use 'public.acme.v1.rims' instead)",
				segment, i+1)
		}
	}

	return nil
}

func ValidateRequestId(requestId string) error {
	if strings.TrimSpace(requestId) != "" && !requestIdRegex.Match([]byte(requestId)) {
		return errors.New("Request ID should be at most 128 characters long and should contain only " +
			"alphanumeric characters, _, space, - or /")
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

// ValidateRimContent validates that the content is either valid JSON or valid JWT (signed/unsigned)
func ValidateRimContent(content string) error {
	if content == "" {
		return errors.New("RIM content cannot be empty")
	}

	// First check if it's valid JSON
	if json.Valid([]byte(content)) {
		return nil
	}

	// If not JSON, check if it's a valid JWT format
	// JWT has format: header.payload.signature (3 parts separated by dots)
	// For unsigned JWT: header.payload. (ends with dot, no signature)
	parts := strings.Split(content, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return errors.New("RIM content must be either valid JSON or valid JWT format")
	}

	// Try to parse as JWT without verification to validate structure
	_, _, err := new(jwt.Parser).ParseUnverified(content, jwt.MapClaims{})
	if err != nil {
		return errors.Wrap(err, "RIM content is not valid JSON or JWT format")
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

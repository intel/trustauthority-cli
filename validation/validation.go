/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package validation

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"intel/amber/tac/v1/constants"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	stringReg = regexp.MustCompile("(^[a-zA-Z0-9_ \\/.-]*$)")
	emailReg  = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+\/=?^_'{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)
	// Regex to validate Amber API key. Key should contain characters between a-z, A-Z, 0-9
	// and should be of size between 30 and 128
	apiKeyRegex = regexp.MustCompile(`^[a-zA-Z0-9]{30,128}$`)
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
	c := filepath.Clean(path)
	r, err := filepath.EvalSymlinks(c)
	if err != nil {
		return c, fmt.Errorf("%s: %s", constants.ErrorInvalidPath, path)
	}
	return r, nil
}

func ValidateSize(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	} else if fi.Size() > constants.MaxPolicyFileSize {
		return fmt.Errorf("%s: %d", constants.ErrorInvalidSize, fi.Size())
	}
	return nil
}

func ValidateAmberAPIKey(apiKey string) error {
	if strings.TrimSpace(apiKey) == "" {
		return errors.Errorf("%s config variable needs to be set with a proper API Key before using CLI", constants.AmberApiKeyEnvVar)
	}
	if matched := apiKeyRegex.MatchString(apiKey); !matched {
		return errors.New("Invalid API key found in configuration file. Please update it with a valid API key.")
	}
	return nil
}

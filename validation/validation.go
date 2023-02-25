/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package validation

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"intel/amber/tac/v1/constants"
	"path/filepath"
	"regexp"
)

var (
	stringReg = regexp.MustCompile("(^[a-zA-Z0-9_ \\/.-]*$)")
	emailReg  = regexp.MustCompile("^[a-zA-Z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}$")
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

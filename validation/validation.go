/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package validation

import (
	"errors"
	"regexp"
)

var (
	stringReg = regexp.MustCompile("(^[a-zA-Z0-9_ \\/.-]*$)")
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

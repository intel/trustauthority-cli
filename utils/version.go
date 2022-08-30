/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 *
 *
 */

package utils

import (
	"encoding/json"
	"intel/amber/tac/v1/constants"
)

var Version = ""
var GitHash = ""
var BuildDate = ""

type CliVersion struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	GitHash   string `json:"gitHash"`
	BuildDate string `json:"buildDate"`
}

var cliVersion = CliVersion{
	Name:      constants.ExplicitCLIName,
	Version:   Version,
	GitHash:   GitHash,
	BuildDate: BuildDate,
}

func GetVersion() (string, error) {
	version, err := json.Marshal(cliVersion)
	if err != nil {
		return "", err
	}

	return string(version), nil
}

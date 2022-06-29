/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package utils

import (
	"bufio"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"intel/amber/tac/v1/constants"
	"io"
	"os"
	"strings"
)

func ReadAnswerFileToEnv(filename string) error {
	fin, err := os.Open(filename)
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

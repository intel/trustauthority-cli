/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package config

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/utils"
	"intel/amber/tac/v1/validation"
	"net/url"
	"os"
	"strings"
)

type Configuration struct {
	AmberBaseUrl      string `yaml:"amber-base-url" mapstructure:"amber-base-url"`
	LogLevel          string `yaml:"log-level" mapstructure:"log-level"`
	HTTPClientTimeout int    `yaml:"http-client-timeout" mapstructure:"http-client-timeout"`
}

// this function sets the configuration file name and type
func init() {
	viper.SetConfigName(constants.ConfigFileName)
	viper.SetConfigType(constants.ConfigFileExtension)
	viper.AddConfigPath(constants.ConfigDir)
}

func LoadConfiguration() (*Configuration, error) {
	ret := Configuration{}
	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			return &ret, errors.Wrap(err, "Config file not found")
		}
		return &ret, errors.Wrap(err, "Failed to load config")
	}
	if err := viper.Unmarshal(&ret); err != nil {
		return &ret, errors.Wrap(err, "Failed to unmarshal config")
	}
	return &ret, nil
}

func (c *Configuration) Save(filename string) error {
	path, err := validation.ValidatePath(filename)
	if err != nil {
		return errors.Wrap(err, "Invalid ConfigurationFilePath")
	}
	configFile, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "Failed to create config file")
	}
	defer func() {
		derr := configFile.Close()
		if derr != nil {
			log.WithError(derr).Error("Error closing config file")
		}
	}()

	err = yaml.NewEncoder(configFile).Encode(c)
	if err != nil {
		return errors.Wrap(err, "Failed to encode config structure")
	}
	return nil
}

func SetupConfig(envFilePath string) error {
	var err error
	if envFilePath == "" {
		return errors.New("EnvFilePath needs to be provided in configuration")
	}
	if _, err = os.Stat(constants.DefaultConfigFilePath); err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(constants.DefaultConfigFilePath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if err = utils.ReadAnswerFileToEnv(envFilePath); err != nil {
		return err
	}

	//set default
	viper.SetDefault(constants.Loglevel, constants.DefaultLogLevel)
	viper.SetDefault(constants.HttpClientTimeout, constants.DefaultHttpClientTimeout)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	configValues := &Configuration{}

	configValues.AmberBaseUrl = viper.GetString(constants.AmberBaseUrl)
	if configValues.AmberBaseUrl == "" {
		return errors.New("Amber base URL needs to be provided in configuration")
	}

	_, err = url.Parse(configValues.AmberBaseUrl)
	if err != nil {
		return errors.Wrap(err, "Invalid Amber Base URL")
	}

	logLevel, err := log.ParseLevel(viper.GetString(constants.Loglevel))
	if err != nil {
		log.Warn("Invalid/No log level provided. Setting log level to info")
	} else {
		configValues.LogLevel = logLevel.String()
	}

	configValues.HTTPClientTimeout = viper.GetInt(constants.HttpClientTimeout)

	if err = configValues.Save(constants.DefaultConfigFilePath); err != nil {
		return err
	}
	return nil
}

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
	"intel/tac/v1/constants"
	"intel/tac/v1/utils"
	"intel/tac/v1/validation"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Configuration struct {
	TrustAuthorityBaseUrl string `yaml:"trustauthority-url" mapstructure:"trustauthority-url"`
	TrustAuthorityApiKey  string `yaml:"trustauthority-api-key" mapstructure:"trustauthority-api-key"`
	LogLevel              string `yaml:"log-level" mapstructure:"log-level"`
	HTTPClientTimeout     int    `yaml:"http-client-timeout" mapstructure:"http-client-timeout"`
}

// this function sets the configuration file name and type
func init() {
	userHomeDir, _ := os.UserHomeDir()
	viper.SetConfigName(constants.ConfigFileName)
	viper.SetConfigType(constants.ConfigFileExtension)
	viper.AddConfigPath(userHomeDir + constants.ConfigDir)
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

func SetupConfig(envFilePath string) error {
	if envFilePath == "" {
		return errors.New("EnvFilePath needs to be provided in configuration")
	}

	_, err := validation.ValidatePath(envFilePath)
	if err != nil {
		return errors.Wrap(err, "Invalid Env file path provided")
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "Error fetching user home directory path")
	}

	cleanedConfigPath := filepath.Clean(userHomeDir + constants.DefaultConfigFilePath)

	configFile, err := os.OpenFile(cleanedConfigPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, constants.DefaultFilePermission)
	if err != nil {
		return errors.Wrap(err, "Failed to open/create config file")
	}
	defer func() {
		derr := configFile.Close()
		if derr != nil {
			log.WithError(derr).Error("Error closing config file")
		}
	}()

	if err = utils.ReadAnswerFileToEnv(envFilePath); err != nil {
		return err
	}

	//set default
	viper.SetDefault(constants.Loglevel, constants.DefaultLogLevel)
	viper.SetDefault(constants.HttpClientTimeout, constants.DefaultHttpClientTimeout)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	configValues := &Configuration{}

	configValues.TrustAuthorityBaseUrl = viper.GetString(constants.TrustAuthBaseUrl)
	if configValues.TrustAuthorityBaseUrl == "" {
		return errors.New("Trust Authority base URL needs to be provided in configuration")
	}

	if !strings.HasPrefix(configValues.TrustAuthorityBaseUrl, "https://") {
		return errors.New("Invalid base URL, must start with 'https://'")
	}

	_, err = url.Parse(configValues.TrustAuthorityBaseUrl)
	if err != nil {
		return errors.Wrap(err, "Invalid Trust Authority Base URL")
	}

	configValues.TrustAuthorityApiKey = viper.GetString(constants.TrustAuthApiKeyEnvVar)
	if err := validation.ValidateTrustAuthorityAPIKey(configValues.TrustAuthorityApiKey); err != nil {
		return errors.Wrap(err, "Invalid API Key provided")
	}

	logLevel, err := log.ParseLevel(viper.GetString(constants.Loglevel))
	if err != nil {
		log.Warn("Invalid/No log level provided. Setting log level to info")
	} else {
		configValues.LogLevel = logLevel.String()
	}

	configValues.HTTPClientTimeout = viper.GetInt(constants.HttpClientTimeout)

	err = yaml.NewEncoder(configFile).Encode(configValues)
	if err != nil {
		return errors.Wrap(err, "Failed to encode config structure")
	}
	return nil
}

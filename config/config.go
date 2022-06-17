package config

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"intel/amber/tac/v1/constants"
	"os"
)

type Configuration struct {
	AmberBaseUrl      string    `yaml:"amber-base-url" mapstructure:"amber-base-url"`
	TenantId          string    `yaml:"tenant-id" mapstructure:"tenant-id"`
	LogLevel          log.Level `yaml:"log"`
	HTTPClientTimeout int       `yaml:"http-client-timeout" mapstructure:"http-client-timeout"`
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
	configFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
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

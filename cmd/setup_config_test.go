package cmd

import (
	"github.com/stretchr/testify/assert"
	"intel/tac/v1/config"
	"intel/tac/v1/utils"
	"intel/tac/v1/validation"
	"testing"
)

func TestSetupConfigCmdEmptyEnvFile(t *testing.T) {
	//Empty Env file path provided
	err := setupConfigCmd.RunE(nil, []string{""})
	if err != nil {
		assert.Error(t, err)
	}
}

func TestSetupConfigInvalidEnvFileError(t *testing.T) {
	//Invalid EnvFilePath provided
	err := config.SetupConfig("test")
	if err != nil {
		assert.Error(t, err)
	}
}

func TestSetupConfigInValidEnvFile(t *testing.T) {
	//Invalid linux file path provided
	err := config.SetupConfig("../test/resou@ces/env_file.env")
	if err != nil {
		assert.Error(t, err)
	}
}

func TestSetupConfigReadAnswerFileToEnvFileError(t *testing.T) {
	//Invalid ReadAnswerFileToEnvPath provided
	err := utils.ReadAnswerFileToEnv("test")
	if err != nil {
		assert.Error(t, err)
	}
}

func TestSetupConfigInvalidURL(t *testing.T) {
	//Invalid Trust authority URL provided
	err := validation.ValidateURL("http://amber-dev02-user12.project-amber-smas.com/")
	if err != nil {
		assert.Error(t, err)

	}
}

func TestSetupConfigErrorParsingURL(t *testing.T) {
	//Error parsing Trust authority URL
	err := validation.ValidateURL("Segment%%2815197306101420000%29.ts")
	if err != nil {
		assert.Error(t, err)

	}
}

func TestSetupConfigvalidURL(t *testing.T) {
	//Valid Trust authority URL provided
	err := validation.ValidateURL("https://amber-dev02-user12.project-amber-smas.com/")
	if err != nil {
		assert.NoError(t, err)

	}
}

func TestSetupConfigInvalidAPIKey(t *testing.T) {
	//Invalid API key provided
	err := validation.ValidateTrustAuthorityAPIKey("test")
	if err != nil {
		assert.Error(t, err)
	}

}

func TestSetupConfigEmptyAPIKey(t *testing.T) {
	//Empty API key provided
	err := validation.ValidateTrustAuthorityAPIKey("")
	if err != nil {
		assert.Error(t, err)
	}

}

func TestSetupConfigValidAPIKey(t *testing.T) {
	//Valid API key provided
	err := validation.ValidateTrustAuthorityAPIKey("djE6MDM1ZmI4NjgtYThjOC00NTJkLTg4ZjYtNmFjMWM5MWJkODI0OmM2czFaeldxOGI5VmEwRXRrbFNMbzJLY0gwM0xtbEFtUnZSQzIyaDE")
	if err != nil {
		assert.NoError(t, err)
	}

}

func TestSetupConfigEmptyClientAPIName(t *testing.T) {
	//Empty API key provided
	err := validation.ValidateApiClientName("")
	if err != nil {
		assert.Error(t, err)
	}

}

func TestSetupConfigEmptyTagName(t *testing.T) {
	// Empty Tag Name provided
	err := validation.ValidateTagName("")
	if err != nil {
		assert.Error(t, err)
	}

}

func TestSetupConfigInvalidTagValue(t *testing.T) {
	// Invalid Tag Value provided
	err := validation.ValidateTagValue("test@@@")
	if err != nil {
		assert.Error(t, err)
	}

}

func TestSetupConfigEmptyTagValue(t *testing.T) {
	//Empty Tag Value provided
	err := validation.ValidateTagValue("")
	if err != nil {
		assert.Error(t, err)
	}
}

func TestSetupConfigEmptyPolicyName(t *testing.T) {
	//Empty Policy Name provided
	err := validation.ValidatePolicyName("")
	if err != nil {
		assert.Error(t, err)
	}

}

func TestSetupConfigInvalidFileSize(t *testing.T) {
	//Invalid File Size provided
	err := validation.ValidateSize("test")
	if err != nil {
		assert.Error(t, err)
	}
}

func TestCheckKeyFileInvalidPrivateKeyPath(t *testing.T) {
	//Invalid Private Key Path provided
	_, _, err := utils.CheckKeyFiles("/.privateKey.txt", "./certKey.txt")
	if err != nil {
		assert.Error(t, err, "")
	}

}

func TestInvalidOutputGeneration(t *testing.T) {
	//Invalid Output Generation provided
	_, err := utils.GenerateOutputFileName("/.inputFile")
	if err != nil {
		assert.Error(t, err, "")
	}

}

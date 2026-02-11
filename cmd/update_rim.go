/*
 * Copyright (C) 2026 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"intel/tac/v1/client/pms"
	"intel/tac/v1/config"
	"intel/tac/v1/constants"
	"intel/tac/v1/models"
	"intel/tac/v1/utils"
	"intel/tac/v1/validation"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// updateRimCmd represents the updateRim command
var updateRimCmd = &cobra.Command{
	Use:   constants.RimCmd,
	Short: "Update an existing RIM with JSON or JWT content",
	Long:  `Update an existing Reference Integrity Manifest (RIM) using either plain JSON or signed JWT format`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("update rim called")
		response, err := updateRim(cmd)
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Println("RIM: \n\n", response)
		return nil
	},
}

func init() {
	updateCmd.AddCommand(updateRimCmd)

	updateRimCmd.Flags().StringP(constants.RimIdParamName, "r", "", "ID of the RIM to be updated")
	updateRimCmd.Flags().StringP(constants.RimDescriptionParamName, "d", "", "Description of the RIM (optional)")
	updateRimCmd.Flags().StringP(constants.RimContentFileParamName, "f", "", "Path of the file containing the RIM content in JSON or JWT format. The file size should be <= 20 KB")
	updateRimCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
	updateRimCmd.MarkFlagRequired(constants.RimIdParamName)
}

func updateRim(cmd *cobra.Command) (string, error) {
	configValues, err := config.LoadConfiguration()
	if err != nil {
		return "", err
	}
	client := &http.Client{
		Timeout: time.Duration(configValues.HTTPClientTimeout) * time.Second,
	}

	pmsUrl, err := url.Parse(configValues.TrustAuthorityBaseUrl + constants.PmsBaseUrl)
	if err != nil {
		return "", err
	}

	if err = setRequestId(cmd); err != nil {
		return "", err
	}

	rimIdString, err := cmd.Flags().GetString(constants.RimIdParamName)
	if err != nil {
		return "", err
	}

	rimId, err := uuid.Parse(rimIdString)
	if err != nil {
		return "", errors.Wrap(err, "Invalid RIM id provided")
	}

	rimDescription, err := cmd.Flags().GetString(constants.RimDescriptionParamName)
	if err != nil {
		return "", err
	}

	rimContentFilePath, err := cmd.Flags().GetString(constants.RimContentFileParamName)
	if err != nil {
		return "", err
	}

	// Validate that at least one of description or content file path is provided
	if rimDescription == "" && rimContentFilePath == "" {
		return "", errors.New("At least one of description or content file path must be provided")
	}

	var rimUpdateReq = models.RimUpdateRequest{
		Id:          rimId,
		Description: rimDescription,
	}

	// Update content if file path provided
	if rimContentFilePath != "" {
		path, err := validation.ValidatePath(rimContentFilePath)
		if err != nil {
			return "", err
		}

		err = validation.ValidateSize(rimContentFilePath)
		if err != nil {
			return "", err
		}
		rimContentBytes, err := os.ReadFile(path)
		if err != nil {
			return "", errors.Wrap(err, "Error reading RIM content file")
		}

		// Trim whitespace
		trimmedContent := strings.TrimSpace(string(rimContentBytes))

		// Validate that content is either valid JSON or valid JWT
		if err := validation.ValidateRimContent(trimmedContent); err != nil {
			return "", err
		}

		// For JWT, wrap as JSON string; for JSON object, use as-is
		var contentForRequest json.RawMessage
		if json.Valid([]byte(trimmedContent)) {
			// Already valid JSON (could be JSON object or JSON string)
			contentForRequest = json.RawMessage(trimmedContent)
		} else {
			// Must be JWT - wrap as JSON string
			jsonString, err := json.Marshal(trimmedContent)
			if err != nil {
				return "", errors.Wrap(err, "Error marshalling RIM content")
			}
			contentForRequest = json.RawMessage(jsonString)
		}

		rimUpdateReq.Content = contentForRequest
	}

	pmsClient := pms.NewPmsClient(client, pmsUrl, apiKey)
	response, err := pmsClient.UpdateRim(&rimUpdateReq)
	if err != nil {
		return "", err
	}

	responseBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

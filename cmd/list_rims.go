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
	"intel/tac/v1/utils"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// getRimsCmd represents the getRims command
var getRimsCmd = &cobra.Command{
	Use:   constants.RimCmd,
	Short: "Get list of RIMs or specific RIM for a tenant",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("list rims called")
		response, err := getRims(cmd)
		utils.PrintRequestAndTraceId()
		if err != nil {
			return err
		}
		fmt.Println("RIMs: \n\n", response)
		return nil
	},
}

func init() {
	listCmd.AddCommand(getRimsCmd)

	getRimsCmd.Flags().StringP(constants.RimIdParamName, "r", "", "ID of the RIM to retrieve")
	getRimsCmd.Flags().StringP(constants.RimNameParamName, "n", "", "Name of the RIM to search for")
	getRimsCmd.Flags().BoolP(constants.IncludeOwnPrivateParamName, "", false, "Include tenant's own private RIMs")
	getRimsCmd.Flags().BoolP(constants.IncludeOwnPublicParamName, "", false, "Include tenant's own public RIMs")
	getRimsCmd.Flags().BoolP(constants.IncludeOthersPublicParamName, "", false, "Include others' public RIMs")
	getRimsCmd.Flags().BoolP(constants.IncludeOnlyOthersPublicParamName, "", false, "Include only others' public RIMs")
	getRimsCmd.Flags().BoolP(constants.ExcludeContentParamName, "", false, "Exclude RIM content from response")
	getRimsCmd.Flags().BoolP(constants.ExcludeDescriptionParamName, "", false, "Exclude description from response")
	getRimsCmd.Flags().StringP(constants.RequestIdParamName, "q", "", "Request ID to be associated with the specific request. This is optional.")
}

func getRims(cmd *cobra.Command) (string, error) {

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

	rimName, err := cmd.Flags().GetString(constants.RimNameParamName)
	if err != nil {
		return "", err
	}

	includeOwnPrivate, err := cmd.Flags().GetBool(constants.IncludeOwnPrivateParamName)
	if err != nil {
		return "", err
	}
	includeOwnPublic, err := cmd.Flags().GetBool(constants.IncludeOwnPublicParamName)
	if err != nil {
		return "", err
	}
	includeOthersPublic, err := cmd.Flags().GetBool(constants.IncludeOthersPublicParamName)
	if err != nil {
		return "", err
	}
	includeOnlyOthersPublic, err := cmd.Flags().GetBool(constants.IncludeOnlyOthersPublicParamName)
	if err != nil {
		return "", err
	}
	excludeContent, err := cmd.Flags().GetBool(constants.ExcludeContentParamName)
	if err != nil {
		return "", err
	}
	excludeDescription, err := cmd.Flags().GetBool(constants.ExcludeDescriptionParamName)
	if err != nil {
		return "", err
	}

	pmsClient := pms.NewPmsClient(client, pmsUrl, apiKey)

	var responseBytes []byte
	if rimIdString != "" {
		// Get specific RIM by ID
		rimId, err := uuid.Parse(rimIdString)
		if err != nil {
			return "", errors.Wrap(err, "Invalid RIM id provided")
		}
		response, err := pmsClient.GetRim(rimId)
		if err != nil {
			return "", err
		}

		responseBytes, err = json.MarshalIndent(response, "", "  ")
		if err != nil {
			return "", err
		}
	} else {
		// Search RIMs with query parameters
		queryParams := make(map[string]string)
		if rimName != "" {
			queryParams["rimName"] = rimName
		}
		if includeOwnPrivate {
			queryParams["includeOwnPrivate"] = "true"
		}
		if includeOwnPublic {
			queryParams["includeOwnPublic"] = "true"
		}
		if includeOthersPublic {
			queryParams["includeOthersPublic"] = "true"
		}
		if includeOnlyOthersPublic {
			queryParams["includeOnlyOthersPublic"] = "true"
		}
		if excludeContent {
			queryParams["excludeContent"] = "true"
		}
		if excludeDescription {
			queryParams["excludeDescription"] = "true"
		}

		response, err := pmsClient.SearchRim(queryParams)
		if err != nil {
			return "", err
		}
		responseBytes, err = json.MarshalIndent(response, "", "  ")
		if err != nil {
			return "", err
		}
	}

	return string(responseBytes), nil
}

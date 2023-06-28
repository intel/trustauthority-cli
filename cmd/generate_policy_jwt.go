/*
 * Copyright (C) 2023 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"fmt"
	"github.com/fatih/set"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"intel/amber/tac/v1/constants"
	"intel/amber/tac/v1/models"
	"intel/amber/tac/v1/utils"
	"intel/amber/tac/v1/validation"
	"os"
)

// createPolicyJwtCmd represents the createPolicyJwtCmd command
var createPolicyJwtCmd = &cobra.Command{
	Use:   constants.PolicyJwtCmd,
	Short: "Generates signed/unsigned JWT for the Rego policy provided",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("create policy-jwt called")
		err := generatePolicyJwt(cmd)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	createCmd.AddCommand(createPolicyJwtCmd)

	createPolicyJwtCmd.Flags().StringP(constants.PolicyFileParamName, "f", "", "Path of the file containing the rego policy to be uploaded")
	createPolicyJwtCmd.Flags().BoolP(constants.SignObjectParamName, "s", false, "Determines if the JWT needs to be signed. Generates a JWS when this parameter is set")
	createPolicyJwtCmd.Flags().StringP(constants.PrivateKeyFileParamName, "p", "", "Path of the file containing the private key to be used to sign the policy. To be used only if -s (sign) parameter is set, else it is ignored")
	createPolicyJwtCmd.Flags().StringP(constants.CertificateFileParamName, "c", "", "Path of the file containing the certificate to be added to the JWT. To be used only if -s (sign) parameter is set, else it is ignored")
	createPolicyJwtCmd.Flags().StringP(constants.AlgorithmParamName, "a", constants.PS384, "Algorithm to be used to sign Amber JWT policy (RS256|PS256|RS384|PS384). To be used only if -s (sign) parameter is set, else it is ignored")
	createPolicyJwtCmd.MarkFlagRequired(constants.PolicyFileParamName)
}

func generatePolicyJwt(cmd *cobra.Command) error {
	var tokenString, algorithm string
	// Create permitted algorithm set
	algorithms := set.New(set.NonThreadSafe)
	policyFilePath, err := cmd.Flags().GetString(constants.PolicyFileParamName)
	if err != nil {
		return err
	}
	if policyFilePath == "" {
		return errors.New("Policy file path cannot be empty")
	}

	path, err := validation.ValidatePath(policyFilePath)
	if err != nil {
		return errors.Wrap(err, "Invalid policy file path provided")
	}

	policyBytes, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "Error reading policy file")
	}

	if len(policyBytes) == 0 {
		return errors.New("Policy file does not contain a rego policy")
	}
	claims := models.PolicyClaims{
		AttestationPolicy: string(policyBytes),
	}

	algorithms.Add(constants.RS256, constants.PS256, constants.RS384, constants.PS384)

	signJwt, err := cmd.Flags().GetBool(constants.SignObjectParamName)
	if err != nil {
		return err
	}

	if signJwt {
		algorithm, err = cmd.Flags().GetString(constants.AlgorithmParamName)
		if err != nil {
			return err
		}
		if !algorithms.Has(algorithm) {
			return errors.New("Input algorithm is not supported")
		}

		privateKeyFilePath, err := cmd.Flags().GetString(constants.PrivateKeyFileParamName)
		if err != nil {
			return err
		}

		certFilePath, err := cmd.Flags().GetString(constants.CertificateFileParamName)
		if err != nil {
			return err
		}

		privKeyFinal, certContents, err := utils.CheckKeyFiles(privateKeyFilePath, certFilePath)
		if err != nil {
			return err
		}

		// Check if provided algorithm makes sense
		signMethod := utils.CheckSigningAlgorithm(privKeyFinal, algorithm)
		if signMethod == nil {
			return errors.New("Signing algorithm provided as input is not compatible with the private key type")
		}

		signedToken := &jwt.Token{
			Header: map[string]interface{}{
				"alg": signMethod.Alg(),
			},
			Claims: claims,
			Method: signMethod,
		}
		signedToken.Header[constants.KeyHeader] = []string{certContents}
		tokenString, err = signedToken.SignedString(privKeyFinal)
		if err != nil {
			return err
		}
	} else {
		algorithm = constants.NonAlg
		token := &jwt.Token{
			Header: map[string]interface{}{
				"alg": jwt.SigningMethodNone.Alg(),
			},
			Claims: claims,
			Method: jwt.SigningMethodNone,
		}

		tokenString, err = token.SigningString()
		if err != nil {
			return err
		}
		tokenString = tokenString + "."
	}

	outputFile, err := utils.GenerateOutputFileName(policyFilePath)
	if err != nil {
		return err
	}
	// Output to screen
	generateConsoleOutput(string(policyBytes), algorithm, outputFile, tokenString)

	// Output to file
	err = os.WriteFile(outputFile, []byte(tokenString), 0400)
	if err != nil {
		return err
	}
	return nil
}

// Print out the contents on console
func generateConsoleOutput(policy, algorithm, outputFile, policyToken string) {
	fmt.Println("Original policy:")
	fmt.Println(policy)
	fmt.Println("Algorithm used during signing: ", algorithm)
	if outputFile != "" {
		fmt.Println("Policy token is stored in file ", outputFile)
	}
	fmt.Println("Policy token generated:")
	fmt.Println(policyToken)
}

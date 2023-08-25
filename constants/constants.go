/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package constants

import (
	"errors"
)

const (
	BinDir                = "/.local/bin/trustauthorityctl"
	ConfigDir             = "/.config/trustauthorityctl/"
	LogDir                = ConfigDir + "logs/"
	DefaultConfigFilePath = ConfigDir + "config.yaml"
	ConfigFileName        = "config"
	ConfigFileExtension   = "yaml"
	LogFilePath           = LogDir + "trustauthorityctl.log"
	DefaultFilePermission = 0640
	MaxPolicyFileSize     = 10240
	ExplicitCLIName       = "Intel Trust Authority CLI"
)

// Command and parameter names
const (
	RequestIdParamName           = "request-id"
	UserIdParamName              = "user-id"
	ServiceIdParamName           = "service-id"
	ServiceOfferIdParamName      = "service-offer-id"
	ProductIdParamName           = "product-id"
	PlanIdParamName              = "plan-id"
	ApiClientIdParamName         = "api-client-id"
	TagIdParamName               = "tag-id"
	ActivationStatus             = "status"
	PolicyIdsParamName           = "policy-ids"
	ApiClientNameParamName       = "api-client-name"
	EmailIdParamName             = "email-id"
	UserRoleParamName            = "user-role"
	PolicyFileParamName          = "policy-file"
	PrivateKeyFileParamName      = "privkeyfile"
	CertificateFileParamName     = "certfile"
	SignObjectParamName          = "sign"
	PolicyIdParamName            = "policy-id"
	PolicyNameParamName          = "policy-name"
	PolicyTypeParamName          = "policy-type"
	AttestationTypeParamName     = "attestation-type"
	TagNameParamName             = "tag-name"
	TagKeyAndValuesParamName     = "tag-key-value"
	EnvFileParamName             = "env-file"
	AlgorithmParamName           = "algorithm"
	DisableNotificationParamName = "disable-notification"

	RootCmd        = "trustauthorityctl"
	CreateCmd      = "create"
	ListCmd        = "list"
	DeleteCmd      = "delete"
	UpdateCmd      = "update"
	UninstallCmd   = "uninstall"
	VersionCmd     = "version"
	SetupConfigCmd = "config"
)

// Resource names
const (
	PolicyCmd         = "policy"
	PolicyJwtCmd      = "policy-jwt"
	UserCmd           = "user"
	ProductCmd        = "product"
	ServiceOfferCmd   = "serviceOffer"
	ServiceCmd        = "service"
	PlanCmd           = "plan"
	ApiClientCmd      = "apiClient"
	TagCmd            = "tag"
	RoleCmd           = "role"
	TenantSettingsCmd = "tenant-settings"
)

const (
	TrustAuthBaseUrl      = "trustauthority-url"
	TrustAuthApiKeyEnvVar = "trustauthority-api-key"
	HttpClientTimeout     = "http-client-timeout"
	Loglevel              = "log-level"

	DefaultLogLevel          = "info"
	DefaultHttpClientTimeout = 10
	DefaultRetryWaitMin      = 2  //minimum time to wait before retry
	DefaultRetryWaitMax      = 10 //maximum time to wait before retry
	DefaultRetryCount        = 2  // number of retries
	ApiClientStatusActive    = "Active"
	ApiClientStatusInactive  = "Inactive"
	ApiClientStatusCancelled = "Cancelled"
	TenantAdminRole          = "Tenant Admin"
	UserRole                 = "User"

	PS384       = "PS384"
	RS256       = "RS256"
	PS256       = "PS256"
	RS384       = "RS384"
	PublicKey   = "PUBLIC KEY"
	CertType    = "CERTIFICATE"
	HashSize256 = "256"
	HashSize384 = "384"
	NonAlg      = "None"
	KeyHeader   = "x5c"
	TimeLayout  = "20060102150405"
)

// HTTP constants
const (
	HTTPMediaTypeJson        = "application/json"
	HTTPHeaderKeyContentType = "Content-Type"
	HTTPHeaderKeyAccept      = "Accept"
	HTTPHeaderKeyApiKey      = "x-api-key"
	HTTPHeaderKeyRequestId   = "request-id"
	HTTPHeaderKeyTraceId     = "trace-id"
)

// API endpoint
const (
	TmsBaseUrl                = "/management/v1"
	PmsBaseUrl                = "/management/v1"
	PolicyApiEndpoint         = "/policies"
	ServiceApiEndpoint        = "/services"
	ServiceOfferApiEndpoint   = "/service-offers"
	ApiClientResourceEndpoint = "/api-clients"
	UserApiEndpoint           = "/users"
	ProductApiEndpoint        = "/products"
	TagApiEndpoint            = "/tags"
	PlanApiEndpoint           = "/plans"
	TenantsApiEndpoint        = "/tenants"
	SettingsEndpoint          = "/settings"
)

type ProductType string

const (
	Attestation ProductType = "attestation"
	Management  ProductType = "management"

	SgxAttestationType           = "SGX Attestation"
	TdxAttestationType           = "TDX Attestation"
	AppraisalPolicyType          = "Appraisal policy"
	TokenCustomizationPolicyType = "Token customization policy"
)

var (
	ErrorInvalidPath        = errors.New("Unsafe or invalid path specified")
	ErrorInvalidSize        = errors.New("Policy File size is  greater than allowed size")
	ServiceUnavailableError = `service unavailable`
)

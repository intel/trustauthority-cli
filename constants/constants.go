/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package constants

const (
	ExecLink              = "/usr/bin/tenantctl"
	HomeDir               = "/opt/tac/"
	ConfigDir             = "/etc/tac/"
	LogDir                = "/var/log/tac/"
	DefaultConfigFilePath = ConfigDir + "config.yaml"
	ConfigFileName        = "config"
	ConfigFileExtension   = "yaml"
	LogFilePath           = LogDir + "tac.log"
	DefaultFilePermission = 0640
	ExplicitCLIName       = "Tenant CLI"
)

//Command and parameter names
const (
	ApiKeyParamName                  = "api-key"
	TenantIdParamName                = "tenant-id"
	UserIdParamName                  = "user-id"
	ServiceIdParamName               = "service-id"
	ServiceOfferIdParamName          = "service-offer-id"
	ProductIdParamName               = "product-id"
	SubscriptionIdParamName          = "subscription-id"
	ActivationStatus                 = "status"
	PolicyIdsParamName               = "policy-ids"
	SubscriptionDescriptionParamName = "subscription-description"
	ServiceNameParamName             = "service-name"
	EmailIdParamName                 = "email-id"
	UserRoleParamName                = "user-role"
	PolicyFileParamName              = "policy-file"
	PolicyIdParamName                = "policy-id"
	TagNameParamName                 = "tag-name"
	TagIdAndValuesParamName          = "tag-id-value"
	EnvFileParamName                 = "env-file"

	RootCmd        = "tenantctl"
	CreateCmd      = "create"
	ListCmd        = "list"
	DeleteCmd      = "delete"
	UpdateCmd      = "update"
	UninstallCmd   = "uninstall"
	VersionCmd     = "version"
	SetupConfigCmd = "config"
)

//Resource names
const (
	PolicyCmd       = "policy"
	UserCmd         = "user"
	ProductCmd      = "product"
	ServiceOfferCmd = "serviceOffer"
	ServiceCmd      = "service"
	SubscriptionCmd = "subscription"
	TagCmd          = "tag"
	RoleCmd         = "role"
)

const (
	AmberBaseUrl      = "amber-base-url"
	TenantId          = "tenant-id"
	HttpClientTimeout = "http-client-timeout"
	Loglevel          = "log-level"

	DefaultLogLevel          = "info"
	DefaultHttpClientTimeout = 10
)

//HTTP constants
const (
	HTTPMediaTypeJson        = "application/json"
	HTTPHeaderKeyTenantId    = "TenantId"
	HTTPHeaderKeyContentType = "Content-Type"
	HTTPHeaderKeyAccept      = "Accept"
	HTTPHeaderKeyApiKey      = "X-API-KEY"
	HTTPHeaderKeyCreatedBy   = "Created-By"
	HTTPHeaderKeyUpdatedBy   = "Updated-By"
)

//API endpoint
const (
	TmsBaseUrl              = "/tms/v1"
	PmsBaseUrl              = "/ps/v1"
	PolicyApiEndpoint       = "/policies"
	TenantApiEndpoint       = "/tenants"
	ServiceApiEndpoint      = "/services"
	ServiceOfferApiEndpoint = "/serviceOffers"
	SubscriptionApiEndpoint = "/subscriptions"
	UserApiEndpoint         = "/users"
	ProductApiEndpoint      = "/products"
	TagApiEndpoint          = "/tags"
	TagsValuesEndpoint      = "/tags-values"
)

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
)

//Command and parameter names
const (
	ApiKeyParamName             = "api-key"
	TenantIdParamName           = "tenant-id"
	UserIdParamName             = "user-id"
	ServiceIdParamName          = "service-id"
	ServiceOfferIdParamName     = "service-offer-id"
	ProductIdParamName          = "product-id"
	SubscriptionIdParamName     = "subscription-id"
	SubscriptionParamName       = "subscription"
	ServiceDescriptionParamName = "service-description"
	EmailIdParamName            = "email-id"
	UserRoleParamName           = "user-role"
	PolicyFileParamName         = "policy-file"
	PolicyIdParamName           = "policy-id"
	EnvFileParamName            = "env-file"

	RootCmd        = "tenantctl"
	CreateCmd      = "create"
	ListCmd        = "list"
	DeleteCmd      = "delete"
	UpdateCmd      = "update"
	UninstallCmd   = "uninstall"
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
)

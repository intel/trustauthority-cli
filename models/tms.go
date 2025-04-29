/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package models

import (
	"github.com/google/uuid"
	"intel/tac/v1/constants"
	"time"
)

type ApiClientStatus string

type (
	//Tenant - Tenant details response payload
	Tenant struct {
		ID               uuid.UUID         `json:"id"`
		Name             string            `json:"name"`
		Company          string            `json:"company"`
		Address          string            `json:"address"`
		ParentId         uuid.UUID         `json:"-"`
		Email            string            `json:"email"`
		Type             string            `json:"type"`
		TenantAttributes map[string]string `json:"tenant_attributes" gorm:"-"`
		ExternalId       uuid.UUID         `json:"-"`
		SourceId         uuid.UUID         `json:"-"`
		UpdatedAt        time.Time         `json:"-"`
		UpdatedBy        uuid.UUID         `json:"-"`
	}

	Role struct {
		ID             uuid.UUID                 `json:"id"`
		Name           string                    `json:"name"`
		Permission     map[string]PermissionAttr `json:"permissions,omitempty"`
		Scope          string                    `json:"scope,omitempty"`
		ServiceOfferId *uuid.UUID                `json:"service_offer_id,omitempty"`
	}

	PermissionAttr struct {
		Grants []string            `json:"grants"`
		Data   map[string][]string `json:"data,omitempty"`
	}

	TenantRoles struct {
		TenantId   uuid.UUID `json:"tenant_id"`
		SourceId   uuid.UUID `json:"source_id"`
		SourceName string    `json:"source_name"`
		Roles      []Role    `json:"roles"`
	}

	// User details response payload
	User struct {
		ID                     uuid.UUID     `json:"id"`
		Email                  string        `json:"email"`
		RepOf                  []uuid.UUID   `json:"rep_of,omitempty"`
		TenantRoles            []TenantRoles `json:"tenant_roles"`
		Active                 bool          `json:"active"`
		CreatedAt              time.Time     `json:"created_at"`
		PrivacyAcknowledgement bool          `json:"privacy_acknowledgement"`
		CreatorType            string        `json:"-"`
		OldRole                string        `json:"-"`
		OemName                string        `json:"oem_name,omitempty"`
	}

	// TenantUser - Tenant user details response payload
	TenantUser struct {
		ID                     uuid.UUID `json:"id"`
		Email                  string    `json:"email"`
		Role                   Role      `json:"role"`
		Active                 bool      `json:"active"`
		CreatedAt              time.Time `json:"created_at"`
		PrivacyAcknowledgement bool      `json:"privacy_acknowledgement"`
		CreatorType            string    `json:"-"`
		Token                  string    `json:"token"`
	}

	UserRole struct {
		RoleId   uuid.UUID `json:"role_id"`
		UserId   uuid.UUID `json:"user_id"`
		TenantId uuid.UUID `json:"tenant_id"`
	}

	CreateTenantUser struct {
		ID                     uuid.UUID `json:"-"`
		Email                  string    `json:"email"`
		FirstName              string    `json:"firstName"`
		LastName               string    `json:"lastName"`
		TenantId               uuid.UUID `json:"-"`
		UserId                 uuid.UUID `json:"-"`
		SubscriptionId         uuid.UUID `json:"-"`
		Active                 bool      `json:"-"`
		PrivacyAcknowledgement bool      `json:"-"`
		CreatedBy              uuid.UUID `json:"-"`
		CreatorType            string    `json:"-"`
		ClaimsRole             string    `json:"-"`
		Role                   string    `json:"role"`
	}

	TagCreate struct {
		ID             uuid.UUID `json:"-"`
		Name           string    `json:"name"`
		TenantId       uuid.UUID `json:"-"`
		UserId         uuid.UUID `json:"-"`
		SubscriptionId uuid.UUID `json:"-"`
		Predefined     bool      `json:"-"`
		CreatedBy      uuid.UUID `json:"-"`
		CreatorType    string    `json:"-"`
	}

	Tag struct {
		ID          *uuid.UUID `json:"id,omitempty"`
		Name        string     `json:"name"`
		TenantId    uuid.UUID  `json:"-"`
		Predefined  bool       `json:"predefined"`
		CreatedBy   uuid.UUID  `json:"-"`
		CreatorType string     `json:"-"`
	}

	UpdateTenantUserRoles struct {
		UserId            uuid.UUID   `json:"-"`
		TenantId          uuid.UUID   `json:"-"`
		JwtSubscriptionId uuid.UUID   `json:"-"`
		JwtUserId         uuid.UUID   `json:"-"`
		Role              string      `json:"role"`
		RoleIds           []uuid.UUID `json:"-"`
		UpdatedBy         uuid.UUID   `json:"-"`
		UpdaterType       string      `json:"-"`
		AccountManagerId  uuid.UUID   `json:"-"`
		CliamsRole        string      `json:"-"`
	}

	Product struct {
		ID             uuid.UUID             `json:"id"`
		ServiceOfferId uuid.UUID             `json:"service_offer_id"`
		Name           string                `json:"name"`
		Policy         *ProductPolicy        `json:"policy"`
		ExternalId     string                `json:"-"`
		PlanIds        []uuid.UUID           `json:"plan_ids"`
		ProductType    constants.ProductType `json:"product_type"`
	}

	ProductPolicy struct {
		Limit              int `json:"limit"`
		Quota              int `json:"quota"`
		LimitRenewalInSecs int `json:"limit_renewal_period"`
		QuotaRenewalInSecs int `json:"quota_renewal_period"`
	}

	// ServiceOffer - Service Offer details response payload
	ServiceOffer struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	// Service - Service details response payload
	Service struct {
		ID                       uuid.UUID              `json:"id"`
		TenantId                 uuid.UUID              `json:"tenant_id"`
		ServiceOfferId           uuid.UUID              `json:"service_offer_id"`
		Name                     string                 `json:"name"`
		PlanId                   uuid.UUID              `json:"plan_id"`
		PlanName                 string                 `json:"plan_name"`
		Active                   bool                   `json:"active"`
		ExternalId               uuid.UUID              `json:"-"`
		CreatedAt                time.Time              `json:"created_at"`
		CreatorType              string                 `json:"-"`
		ServiceOfferPlanSourceId uuid.UUID              `json:"-"`
		DurationMonths           int                    `json:"-"`
		Attributes               map[string]interface{} `json:"attributes"`
	}

	ServiceDetail struct {
		ID                       uuid.UUID `json:"id"`
		ServiceOfferId           uuid.UUID `json:"service_offer_id"`
		ServiceOfferName         string    `json:"service_offer_name"`
		Name                     string    `json:"name"`
		CreatedAt                time.Time `json:"created_at"`
		Active                   bool      `json:"active"`
		PlanId                   uuid.UUID `json:"plan_id"`
		PlanName                 string    `json:"plan_name"`
		ServiceOfferPlanSourceId uuid.UUID `json:"-"`
	}

	// ApiClient - API Client details response payload
	ApiClient struct {
		ID                uuid.UUID             `json:"id"`
		ServiceId         uuid.UUID             `json:"service_id"`
		ProductId         uuid.UUID             `json:"product_id"`
		ProductName       string                `json:"product_name"`
		Status            ApiClientStatus       `json:"status"`
		Name              string                `json:"name"`
		CreatedAt         time.Time             `json:"created_at"`
		ProductExternalId string                `json:"-"`
		ProductType       constants.ProductType `json:"product_type"`
	}

	// UpdateApiClient - API Client details request to be updated
	UpdateApiClient struct {
		Id             uuid.UUID               `json:"-"`
		ProductId      uuid.UUID               `json:"product_id"`
		ServiceId      uuid.UUID               `json:"-"`
		Name           *string                 `json:"name"`
		TenantId       uuid.UUID               `json:"-"`
		UserId         uuid.UUID               `json:"-"`
		SubscriptionId uuid.UUID               `json:"-"`
		PolicyIds      []uuid.UUID             `json:"policy_ids"`
		TagIdsValues   []ApiClientTagIdValue   `json:"tags"`
		UpdatedBy      uuid.UUID               `json:"-"`
		Status         *ApiClientStatus        `json:"status"`
		UpdaterType    string                  `json:"-"`
		ProductType    []constants.ProductType `json:"-"`
		ExternalId     *string                 `json:"-"`
		UpdateDeleted  bool                    `json:"-"`
	}

	// ApiClientDetail - API Client details response payload
	ApiClientDetail struct {
		ID                uuid.UUID             `json:"id"`
		ServiceId         uuid.UUID             `json:"service_id"`
		ServiceOfferName  string                `json:"service_offer_name"`
		ProductId         uuid.UUID             `json:"product_id"`
		ProductName       string                `json:"product_name"`
		Status            ApiClientStatus       `json:"status"`
		Name              string                `json:"name"`
		Keys              []string              `json:"keys"`
		PolicyIds         []uuid.UUID           `json:"policy_ids"`
		TagsValues        []ApiClientTagValue   `json:"tags"`
		CreatedAt         time.Time             `json:"created_at"`
		ProductExternalId string                `json:"-"`
		ProductType       constants.ProductType `json:"product_type"`
	}

	CreateApiClient struct {
		ProductId    uuid.UUID             `json:"product_id"`
		ServiceId    uuid.UUID             `json:"-"`
		TenantId     uuid.UUID             `json:"-"`
		PolicyIds    []uuid.UUID           `json:"policy_ids"`
		TagIdsValues []ApiClientTagIdValue `json:"tags"`
		Name         string                `json:"name"`
		Status       ApiClientStatus       `json:"status"`
		CreatedBy    uuid.UUID             `json:"-"`
		CreatorType  string                `json:"-"`
	}

	ApiClientPolicies struct {
		PolicyIds []uuid.UUID `json:"policy_ids"`
	}

	ApiClientTagValue struct {
		Name       string `json:"key"`
		Value      string `json:"value"`
		Predefined bool   `json:"predefined"`
	}

	ApiClientTags struct {
		TagsValues []ApiClientTagValue `json:"tags"`
	}

	ApiClientTagIdValue struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	Tags struct {
		Tags []Tag `json:"tags"`
	}

	Plan struct {
		ID             uuid.UUID `json:"id"`
		ServiceOfferId uuid.UUID `json:"service_offer_id"`
		Name           string    `json:"name"`
		MaxKey         int       `json:"max_key"`
		MaxTenantAdmin int       `json:"max_tenant_admin"`
		MaxTenantUser  int       `json:"max_tenant_user"`
		MaxPolicy      int       `json:"max_policy"`
		Ledger         bool      `json:"ledger"`
	}

	PlanProducts struct {
		ID             uuid.UUID `json:"id"`
		ServiceOfferId uuid.UUID `json:"service_offer_id"`
		Name           string    `json:"name"`
		MaxKey         int       `json:"max_key"`
		MaxTenantAdmin int       `json:"max_tenant_admin"`
		MaxTenantUser  int       `json:"max_tenant_user"`
		MaxPolicy      int       `json:"max_policy"`
		Ledger         bool      `json:"ledger"`
		Products       []Product `json:"products"`
	}

	AttestationFailureEmail struct {
		// Email address to which Attestation Failure need to be sent
		// example: jane.doe@domain.com
		AttestationFailureEmail string `json:"attest_failure_email"`
	}
)

/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package models

import (
	"github.com/google/uuid"
	"time"
)

type ApiClientStatus string

type (
	Tenant struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		Company  string    `json:"company"`
		Address  string    `json:"address"`
		ParentId uuid.UUID `json:"-"`
		Email    string    `json:"email"`
	}

	Role struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	TenantRoles struct {
		TenantId uuid.UUID `json:"tenant_id"`
		Roles    []Role    `json:"roles"`
	}

	User struct {
		ID                     uuid.UUID     `json:"id"`
		Email                  string        `json:"email"`
		TenantRoles            []TenantRoles `json:"tenant_roles"`
		Active                 bool          `json:"active"`
		CreatedAt              time.Time     `json:"created_at"`
		PrivacyAcknowledgement bool          `json:"privacy_acknowledgement"`
	}

	TenantUser struct {
		ID                     uuid.UUID `json:"id"`
		Email                  string    `json:"email"`
		Role                   Role      `json:"role"`
		Active                 bool      `json:"active"`
		CreatedAt              time.Time `json:"created_at"`
		PrivacyAcknowledgement bool      `json:"privacy_acknowledgement"`
	}

	UserRole struct {
		RoleId   uuid.UUID `json:"role_id"`
		UserId   uuid.UUID `json:"user_id"`
		TenantId uuid.UUID `json:"tenant_id"`
	}

	CreateTenantUser struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		TenantId  uuid.UUID `json:"-"`
		Active    bool      `json:"-"`
		CreatedBy uuid.UUID `json:"-"`
		Role      string    `json:"role"`
	}

	Tag struct {
		ID         *uuid.UUID `json:"id,omitempty"`
		Name       string     `json:"name"`
		TenantId   uuid.UUID  `json:"-"`
		Predefined bool       `json:"predefined"`
		CreatedBy  uuid.UUID  `json:"-"`
	}

	UpdateTenantUserRoles struct {
		UserId    uuid.UUID   `json:"-"`
		TenantId  uuid.UUID   `json:"-"`
		Role      string      `json:"role"`
		RoleIds   []uuid.UUID `json:"-"`
		UpdatedBy uuid.UUID   `json:"-"`
	}

	Product struct {
		ID             uuid.UUID      `json:"id"`
		ServiceOfferId uuid.UUID      `json:"service_offer_id"`
		Name           string         `json:"name"`
		Policy         *ProductPolicy `json:"policy"`
	}

	ProductPolicy struct {
		Limit              int `json:"limit"`
		Quota              int `json:"quota"`
		LimitRenewalInSecs int `json:"limit_renewal_period"`
		QuotaRenewalInSecs int `json:"quota_renewal_period"`
	}

	ServiceOffer struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	Service struct {
		ID             uuid.UUID `json:"id"`
		TenantId       uuid.UUID `json:"tenant_id"`
		ServiceOfferId uuid.UUID `json:"service_offer_id"`
		Name           string    `json:"name"`
		CreatedAt      time.Time `json:"created_at"`
	}

	UpdateService struct {
		Id        uuid.UUID `json:"-"`
		Name      string    `json:"name"`
		TenantId  uuid.UUID `json:"-"`
		UpdatedBy uuid.UUID `json:"-"`
	}

	ServiceDetail struct {
		ID               uuid.UUID `json:"id"`
		ServiceOfferId   uuid.UUID `json:"service_offer_id"`
		ServiceOfferName string    `json:"service_offer_name"`
		Name             string    `json:"name"`
		CreatedAt        time.Time `json:"created_at"`
	}

	CreateService struct {
		ServiceOfferId uuid.UUID `json:"service_offer_id"`
		TenantId       uuid.UUID `json:"-"`
		Name           string    `json:"name"`
		CreatedBy      uuid.UUID `json:"-"`
	}

	ApiClient struct {
		ID          uuid.UUID       `json:"id"`
		ServiceId   uuid.UUID       `json:"service_id"`
		ProductId   uuid.UUID       `json:"product_id"`
		ProductName string          `json:"product_name"`
		Status      ApiClientStatus `json:"status"`
		Name        string          `json:"name"`
		ExpiredAt   time.Time       `json:"expired_at"`
		CreatedAt   time.Time       `json:"created_at"`
	}

	UpdateApiClient struct {
		Id           uuid.UUID             `json:"-"`
		ProductId    uuid.UUID             `json:"product_id"`
		ServiceId    uuid.UUID             `json:"-"`
		Name         string                `json:"name"`
		TenantId     uuid.UUID             `json:"-"`
		PolicyIds    []uuid.UUID           `json:"policy_ids"`
		TagIdsValues []ApiClientTagIdValue `json:"tags"`
		UpdatedBy    uuid.UUID             `json:"-"`
		Status       ApiClientStatus       `json:"status"`
		ExpiredAt    time.Time             `json:"expired_at"`
	}

	ApiClientDetail struct {
		ID               uuid.UUID           `json:"id"`
		ServiceId        uuid.UUID           `json:"service_id"`
		ServiceOfferName string              `json:"service_offer_name"`
		ProductId        uuid.UUID           `json:"product_id"`
		ProductName      string              `json:"product_name"`
		Status           ApiClientStatus     `json:"status"`
		Name             string              `json:"name"`
		ExpiredAt        time.Time           `json:"expired_at"`
		Keys             []string            `json:"keys"`
		PolicyIds        []uuid.UUID         `json:"policy_ids"`
		TagsValues       []ApiClientTagValue `json:"tags"`
		CreatedAt        time.Time           `json:"created_at"`
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
		ExpiredAt    time.Time             `json:"expired_at"`
	}

	ApiClientPolicies struct {
		PolicyIds []uuid.UUID `json:"policy_ids"`
	}

	ApiClientTagValue struct {
		Name       string `json:"key"`
		Value      string `json:"value"`
		Predefined bool   `json:"predefined"`
	}

	ApiClientTagsValues struct {
		TagsValues []ApiClientTagValue `json:"tags"`
	}

	ApiClientTagIdValue struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	Tags struct {
		Tags []Tag `json:"tags"`
	}
)

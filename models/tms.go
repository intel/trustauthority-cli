/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package models

import (
	"github.com/google/uuid"
	"time"
)

type SubscriptionStatus string

const (
	Active    SubscriptionStatus = "Active"
	Inactive  SubscriptionStatus = "Inactive"
	Cancelled SubscriptionStatus = "Cancelled"
)

type (
	Tenant struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		Company  string    `json:"company"`
		Address  string    `json:"address"`
		ParentId uuid.UUID `json:"parent_id,omitempty"`
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
		ID          uuid.UUID     `json:"id"`
		Email       string        `json:"email"`
		TenantRoles []TenantRoles `json:"tenant_roles"`
		Active      bool          `json:"active"`
		CreatedAt   time.Time     `json:"created_at"`
	}

	UserRole struct {
		RoleId   uuid.UUID `json:"role_id"`
		UserId   uuid.UUID `json:"user_id"`
		TenantId uuid.UUID `json:"tenant_id"`
	}

	CreateTenantUser struct {
		Email     string    `json:"email"`
		TenantId  uuid.UUID `json:"-"`
		Active    bool      `json:"-"`
		CreatedBy uuid.UUID `json:"-"`
		Role      string    `json:"role"`
	}

	UpdateTenantUserRoles struct {
		UserId    uuid.UUID   `json:"-"`
		TenantId  uuid.UUID   `json:"-"`
		Roles     []string    `json:"roles"`
		RoleIds   []uuid.UUID `json:"-"`
		CreatedBy uuid.UUID   `json:"-"`
	}

	CreateUser struct {
		Email     string    `json:"email"`
		Active    bool      `json:"-"`
		CreatedBy uuid.UUID `json:"-"`
	}

	UpdateUser struct {
		Id          uuid.UUID `json:"-"`
		Email       string    `json:"email"`
		SetActive   bool      `json:"set_active"`
		SetInactive bool      `json:"set_inactive"`
		UpdatedBy   uuid.UUID `json:"-"`
	}

	DeleteUser struct {
		Id        uuid.UUID `json:"-"`
		TenantId  uuid.UUID `json:"-"`
		UpdatedBy uuid.UUID `json:"-"`
	}

	//ProvisionTenant used to create a tenant and a Tenant Admin in a single step
	ProvisionTenant struct {
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		Company   string    `json:"company"`
		Address   string    `json:"address"`
		Role      string    `json:"-"`
		CreatedBy uuid.UUID `json:"-"`
		ParentId  uuid.UUID `json:"parent_id"`
	}

	Product struct {
		ID             uuid.UUID `json:"id"`
		ServiceOfferId uuid.UUID `json:"service_offer_id"`
		Name           string    `json:"name"`
	}

	CreateProduct struct {
		Name           string    `json:"name"`
		ServiceOfferId uuid.UUID `json:"-"`
		CreatedBy      uuid.UUID `json:"-"`
	}

	ServiceOffer struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	CreateServiceOffer struct {
		Name      string    `json:"name"`
		CreatedBy uuid.UUID `json:"-"`
	}

	Service struct {
		ID             uuid.UUID `json:"id"`
		TenantId       uuid.UUID `json:"tenant_id"`
		ServiceOfferId uuid.UUID `json:"service_offer_id"`
		Description    string    `json:"description"`
	}

	ServiceDetail struct {
		ID               uuid.UUID `json:"id"`
		TenantId         uuid.UUID `json:"tenant_id"`
		TenantName       string    `json:"tenant_name"`
		ServiceOfferId   uuid.UUID `json:"service_offer_id"`
		ServiceOfferName string    `json:"service_offer_name"`
		Description      string    `json:"description"`
	}

	CreateService struct {
		TenantId       uuid.UUID `json:"_"`
		ServiceOfferId uuid.UUID `json:"service_offer_id"`
		Description    string    `json:"description"`
		CreatedBy      uuid.UUID `json:"-"`
	}

	Subscription struct {
		ID          uuid.UUID          `json:"id"`
		ServiceId   uuid.UUID          `json:"service_id"`
		ProductId   uuid.UUID          `json:"product_id"`
		Status      SubscriptionStatus `json:"status"`
		Description string             `json:"description"`
	}

	SubscriptionDetail struct {
		ID          uuid.UUID          `json:"id"`
		ServiceId   uuid.UUID          `json:"service_id"`
		ProductId   uuid.UUID          `json:"product_id"`
		Status      SubscriptionStatus `json:"status"`
		Description string             `json:"description"`
		Keys        []string           `json:"keys"`
	}

	CreateSubscription struct {
		ServiceId   uuid.UUID `json:"_"`
		TenantId    uuid.UUID `json:"_"`
		ProductId   uuid.UUID `json:"product_id"`
		Description string    `json:"description"`
		CreatedBy   uuid.UUID `json:"-"`
	}
)

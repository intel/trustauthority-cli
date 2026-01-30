/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package models

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// PolicyRequest struct defines the policy
type PolicyRequest struct {
	CommonPolicy
}

type PolicyResponse struct {
	CommonPolicy
	CreatorId       *uuid.UUID `json:"creator_id"`
	UpdaterId       *uuid.UUID `json:"updater_id"`
	Deleted         bool       `json:"deleted"`
	CreatedAt       time.Time  `json:"created_time"`
	UpdatedAt       time.Time  `json:"modified_time"`
	PolicyJWT       string     `json:"policy_jwt,omitempty"`
	PolicyHash      string     `json:"policy_hash"`
	PolicySignature string     `json:"policy_signature"`
	Version         string     `json:"version"`
	SignedByTenant  bool       `json:"signed_by_tenant"`
}

type CommonPolicy struct {
	PolicyId        uuid.UUID `json:"policy_id"`
	Policy          string    `json:"policy"`
	TenantId        uuid.UUID `json:"-"`
	PolicyName      string    `json:"policy_name"`
	PolicyType      string    `json:"policy_type"`
	ServiceOfferId  uuid.UUID `json:"service_offer_id"`
	AttestationType string    `json:"attestation_type"`
}

type PolicyUpdateRequest struct {
	PolicyId        uuid.UUID `json:"policy_id"`
	Policy          string    `json:"policy"`
	TenantId        uuid.UUID `json:"-"`
	PolicyName      string    `json:"policy_name"`
	PolicyType      string    `json:"-"`
	AttestationType string    `json:"-"`
	ServiceOfferId  uuid.UUID `json:"-"`
	PolicyJWT       string    `json:"-"`
	UserId          uuid.UUID `json:"-"`
	SubscriptionId  uuid.UUID `json:"-"`
}

type PolicyClaims struct {
	AttestationPolicy string `json:"AttestationPolicy"`
	jwt.RegisteredClaims
}

type CommonRim struct {
	Id           uuid.UUID       `json:"id"`
	Content      json.RawMessage `json:"content,omitempty"`
	TenantId     uuid.UUID       `json:"-"`
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	Public       bool            `json:"public"`
	PublicDomain string          `json:"-"` // Internal field, not exposed in API response
}

type RimCreateRequest struct {
	Id             uuid.UUID       `json:"-"`
	Content        json.RawMessage `json:"content"`
	TenantId       uuid.UUID       `json:"-"`
	Name           string          `json:"name"`
	Description    string          `json:"description,omitempty"`
	Public         bool            `json:"-"` // Determined from name, not from request
	PublicDomain   string          `json:"-"` // Extracted from name if public
	JWT            string          `json:"-"`
	Hash           string          `json:"-"`
	UserId         uuid.UUID       `json:"-"`
	SubscriptionId uuid.UUID       `json:"-"`
}

type RimUpdateRequest struct {
	Id             uuid.UUID       `json:"-"`
	Content        json.RawMessage `json:"content"`
	TenantId       uuid.UUID       `json:"-"`
	Description    string          `json:"description,omitempty"`
	JWT            string          `json:"-"`
	Hash           string          `json:"-"`
	UserId         uuid.UUID       `json:"-"`
	SubscriptionId uuid.UUID       `json:"-"`
}
type SignedRimResponse struct {
	CommonRim
	CreatorId      *uuid.UUID `json:"creator_id"`
	UpdaterId      *uuid.UUID `json:"updater_id"`
	Deleted        bool       `json:"deleted"`
	CreatedAt      time.Time  `json:"created_time"`
	UpdatedAt      time.Time  `json:"modified_time"`
	JWT            string     `json:"jwt,omitempty"`
	Hash           string     `json:"hash,omitempty"`
	Signature      string     `json:"signature,omitempty"`
	Version        string     `json:"version"`
	SignedByTenant bool       `json:"signed_by_tenant"`
}

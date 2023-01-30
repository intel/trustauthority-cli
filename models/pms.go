/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package models

import (
	"github.com/google/uuid"
	"time"
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
}

type CommonPolicy struct {
	PolicyId        uuid.UUID `json:"policy_id"`
	Policy          string    `json:"policy"`
	PolicyName      string    `json:"policy_name"`
	PolicyType      string    `json:"policy_type"`
	ServiceOfferId  uuid.UUID `json:"service_offer_id"`
	AttestationType string    `json:"attestation_type"`
}

type PolicyUpdateRequest struct {
	PolicyId       uuid.UUID `json:"policy_id"`
	Policy         string    `json:"policy"`
	TenantId       uuid.UUID `json:"-"`
	PolicyName     string    `json:"policy_name"`
	UserId         uuid.UUID `json:"_"`
	SubscriptionId uuid.UUID `json:"-"`
}

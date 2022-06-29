/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package test

import (
	"fmt"
	"github.com/gorilla/mux"
	"intel/amber/tac/v1/config"
	"intel/amber/tac/v1/constants"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	policy = `{
    "policy_id": "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
    "policy": "default matches_sgx_policy = false \n\n matches_sgx_policy = true { \n\n quote := input.quote \n quote.isvsvn == 0 \n  isvprodidValues := [0, 2, 3] \n includes_value(isvprodidValues, quote.isvprodid) \n mrsignerValues:= [ \"d412a4f07ef83892a5915fb2ab584be31e186e5a4f95ab5f6950fd4eb8694d7b\"] \n includes_value(mrsignerValues, quote.mrsigner) \n mrenclaveValues := [\"bab91f200038076ac25f87de0ca67472443c2ebe17ed9ba95314e609038f51ab\"] \n includes_value(mrenclaveValues, quote.mrenclave) \n } \n includes_value(policy_values, quote_value) = true { \n\n policy_value := policy_values[x] \n policy_value == quote_value \n } \n",
	"tenant_id": "f04971b7-fb41-4a9e-a06e-4bf6e71f98b3",
	"user_id": "f04971b7-fb41-4a9e-a06e-4bf6e71f98b3",
	"version": "v1",
    "policy_name": "Sample_Policy_SGX",
    "policy_type": "Appraisal",
    "service_offer_name": "SGX",
    "subscription_id": "e8a72b7e-c4b1-4bdc-bf40-68f23c68a2aa"
}`

	policyList = `[{
    "policy_id": "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
    "policy": "default matches_sgx_policy = false \n\n matches_sgx_policy = true { \n\n quote := input.quote \n quote.isvsvn == 0 \n  isvprodidValues := [0, 2, 3] \n includes_value(isvprodidValues, quote.isvprodid) \n mrsignerValues:= [ \"d412a4f07ef83892a5915fb2ab584be31e186e5a4f95ab5f6950fd4eb8694d7b\"] \n includes_value(mrsignerValues, quote.mrsigner) \n mrenclaveValues := [\"bab91f200038076ac25f87de0ca67472443c2ebe17ed9ba95314e609038f51ab\"] \n includes_value(mrenclaveValues, quote.mrenclave) \n } \n includes_value(policy_values, quote_value) = true { \n\n policy_value := policy_values[x] \n policy_value == quote_value \n } \n",
	"tenant_id": "f04971b7-fb41-4a9e-a06e-4bf6e71f98b3",
	"user_id": "f04971b7-fb41-4a9e-a06e-4bf6e71f98b3",
	"version": "v1",
    "policy_name": "Sample_Policy_SGX",
    "policy_type": "Appraisal",
    "service_offer_name": "SGX",
    "subscription_id": "e8a72b7e-c4b1-4bdc-bf40-68f23c68a2aa"
}]`

	user = `{
        "id": "23011406-6f3b-4431-9363-4e1af9af6b13",
        "email": "arijitgh@gmail.com",
        "tenant_roles": [
            {
                "tenant_id": "89120415-6fbc-41c7-b9f2-3b4ba10e87c9",
                "roles": [
                    {
                        "id": "66ec2e33-8cd3-42b1-8963-c7765205446e",
                        "name": "Tenant Admin"
                    }
                ]
            }
        ],
        "active": false,
        "created_at": "2022-06-19T20:02:55.157679Z"
    }`

	userList = `[
    {
        "id": "23011406-6f3b-4431-9363-4e1af9af6b13",
        "email": "arijitgh@gmail.com",
        "tenant_roles": [
            {
                "tenant_id": "89120415-6fbc-41c7-b9f2-3b4ba10e87c9",
                "roles": [
                    {
                        "id": "66ec2e33-8cd3-42b1-8963-c7765205446e",
                        "name": "Tenant Admin"
                    }
                ]
            }
        ],
        "active": false,
        "created_at": "2022-06-19T20:02:55.157679Z"
    }
]`

	serviceOfferList = `[
    {
        "id": "ae3d7720-08ab-421c-b8d4-1725c358f03e",
        "name": "TDX Attestation"
    }
]`

	serviceList = `[
    {
        "id": "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
        "tenant_id": "89120415-6fbc-41c7-b9f2-3b4ba10e87c9",
        "service_offer_id": "ae3d7720-08ab-421c-b8d4-1725c358f03e",
        "description": "Arijit Service"
    }
]`
	service = `{
        "id": "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
        "tenant_id": "89120415-6fbc-41c7-b9f2-3b4ba10e87c9",
        "service_offer_id": "ae3d7720-08ab-421c-b8d4-1725c358f03e",
        "description": "Arijit Service"
    }`

	subscriptionList = `[
    {
        "id": "3780cc39-cce2-4ec2-a47f-03e55b12e259",
        "service_id": "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
        "product_id": "e169d34f-58ce-4717-9b3a-5c66abd33417",
        "status": "",
        "description": "Arijit Subscription"
    }
]`
	subscription = `{
        "id": "3780cc39-cce2-4ec2-a47f-03e55b12e259",
        "service_id": "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
        "product_id": "e169d34f-58ce-4717-9b3a-5c66abd33417",
        "status": "",
        "description": "Arijit Subscription",
		"keys": [
			"9dca50986c414304a4b1ffe202dcf2b0",
			"996a9a6e67814f1784eadb5405bdabf3"
		]
    }`

	productList = `[
    {
        "id": "e169d34f-58ce-4717-9b3a-5c66abd33417",
        "service_offer_id": "ae3d7720-08ab-421c-b8d4-1725c358f03e",
        "name": "Developer"
    }
]`
)

var idReg = fmt.Sprintf("{id:%s}", "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}")

// MockPmsServer for CLI unit testing
func MockPmsServer(t *testing.T) *httptest.Server {
	policyIdExpr := fmt.Sprintf("%s%s", "/ps/v1/policies/", idReg)

	r := mux.NewRouter()

	r.HandleFunc("/ps/v1/policies", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(policy))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPost)

	r.HandleFunc("/ps/v1/policies", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(policyList))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(policyIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(policy))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(policyIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write(nil)
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodDelete)

	r.HandleFunc(policyIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(policy))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPut)

	return httptest.NewServer(r)
}

// MockTmsServer for CLI unit testing
func MockTmsServer(t *testing.T) *httptest.Server {

	userExpr := fmt.Sprintf("%s%s%s", "/tms/v1/tenants/", idReg, "/users")
	userIdExpr := fmt.Sprintf("%s%s%s%s", "/tms/v1/tenants/", idReg, "/users/", idReg)
	serviceExpr := fmt.Sprintf("%s%s%s", "/tms/v1/tenants/", idReg, "/services")
	serviceIdExpr := fmt.Sprintf("%s%s%s%s", "/tms/v1/tenants/", idReg, "/services/", idReg)
	subscriptionExpr := fmt.Sprintf("%s%s%s%s%s", "/tms/v1/tenants/", idReg, "/services/", idReg, "/subscriptions")
	subscriptionIdExpr := fmt.Sprintf("%s%s%s%s%s%s", "/tms/v1/tenants/", idReg, "/services/", idReg, "/subscriptions/", idReg)
	serviceOfferExpr := fmt.Sprintf("/tms/v1/serviceOffers")
	productExpr := fmt.Sprintf("%s%s%s", "/tms/v1/serviceOffers/", idReg, "/products")

	r := mux.NewRouter()
	r.HandleFunc(userExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(user))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPost)

	r.HandleFunc(userExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(userList))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(userIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(user))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(userIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write(nil)
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodDelete)

	r.HandleFunc(subscriptionExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(subscription))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPost)

	r.HandleFunc(subscriptionExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(subscriptionList))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(subscriptionIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(subscription))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(serviceExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(service))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPost)

	r.HandleFunc(serviceExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(serviceList))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(serviceIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(service))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(serviceOfferExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(serviceOfferList))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(productExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(productList))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	return httptest.NewServer(r)
}

//SetupMockConfiguration setting up mock CLI configurations
func SetupMockConfiguration(serverUrl string) *config.Configuration {

	os.Mkdir(constants.ConfigDir, constants.DefaultFilePermission)
	c := &config.Configuration{
		AmberBaseUrl: serverUrl,
		TenantId:     "f04971b7-fb41-4a9e-a06e-4bf6e71f98b3",
	}
	c.Save("config.yaml")

	return c
}

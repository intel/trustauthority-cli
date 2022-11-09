/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package test

import (
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"intel/amber/tac/v1/config"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	policy = `{
    "policy_id": "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
    "policy": "default matches_sgx_policy = false \n\n matches_sgx_policy = true { \n \n input.amber_sgx_isvsvn == 0 \n\n } ",
	"tenant_id": "f04971b7-fb41-4a9e-a06e-4bf6e71f98b3",
	"user_id": "f04971b7-fb41-4a9e-a06e-4bf6e71f98b3",
	"version": "v1",
    "policy_name": "Sample_Policy_SGX",
    "policy_type": "Appraisal policy",
    "service_offer_name": "SGX Attestation",
    "subscription_id": "e8a72b7e-c4b1-4bdc-bf40-68f23c68a2aa" }`

	policyList = `[{
    "policy_id": "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
    "policy": "default matches_sgx_policy = false \n\n matches_sgx_policy = true { \n \n input.amber_sgx_isvsvn == 0 \n\n } ",
	"tenant_id": "f04971b7-fb41-4a9e-a06e-4bf6e71f98b3",
	"user_id": "f04971b7-fb41-4a9e-a06e-4bf6e71f98b3",
	"version": "v1",
    "policy_name": "Sample_Policy_SGX",
    "policy_type": "Appraisal",
    "service_offer_name": "SGX",
    "subscription_id": "e8a72b7e-c4b1-4bdc-bf40-68f23c68a2aa" }]`

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
			"roles": [
			{
				"id": "66ec2e33-8cd3-42b1-8963-c7765205446e",
				"name": "Tenant Admin"
			}],
        "active": false,
        "created_at": "2022-06-19T20:02:55.157679Z"
    }]`

	serviceOfferList = `[
    {
        "id": "ae3d7720-08ab-421c-b8d4-1725c358f03e",
        "name": "TDX Attestation"
    }]`

	serviceList = `[
    {
        "id": "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
        "tenant_id": "89120415-6fbc-41c7-b9f2-3b4ba10e87c9",
        "service_offer_id": "ae3d7720-08ab-421c-b8d4-1725c358f03e",
        "description": "Test Service"
    }]`
	service = `{
        "id": "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
        "tenant_id": "89120415-6fbc-41c7-b9f2-3b4ba10e87c9",
        "service_offer_id": "ae3d7720-08ab-421c-b8d4-1725c358f03e",
        "description": "Test Service"
    }`

	subscriptionList = `[
    {
        "id": "3780cc39-cce2-4ec2-a47f-03e55b12e259",
        "service_id": "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
        "product_id": "e169d34f-58ce-4717-9b3a-5c66abd33417",
        "status": "",
        "description": "Test Subscription"
    }]`

	subscription = `{
        "id": "3780cc39-cce2-4ec2-a47f-03e55b12e259",
        "service_id": "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
        "product_id": "e169d34f-58ce-4717-9b3a-5c66abd33417",
        "status": "",
        "name": "Test Subscription",
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
    }]`

	tag = `{
	    "id": "a9f765e4-0296-4147-be03-1c9deb7c050f",
	    "name": "Frequency",
	    "predefined": false
	}`

	tagList = `{
	    "tags":[
        {
	            "id": "f31aa1bc-99a1-4706-91ff-218e12c49e00",
	            "name": "Workload",
	            "predefined": true
	        },
	        {
	            "id": "c0b2d143-c9f5-4137-88db-1f2a8e666d0c",
	            "name": "Power",
	            "predefined": false
	        }
	    ]
	  }`

	policyIds = `{
	    "policy_ids": [ "c855d8d6-744f-48c6-a06d-a97ef1811a61","5d389ba0-f683-4d42-ad52-c76a94ccabc7" ]
	}`

	tagsValues = `{
	    "tags_values": [
         {
            "id": "f31aa1bc-99a1-4706-91ff-218e12c49e00",
            "name": "Workload",
            "value": "AI"
		 },
         {
            "id": "f31aa1bc-99a1-4706-91ff-218e12c49e00",
            "name": "Workload",
            "value": "EXE Workload"
        }
    ]}`
)

var idReg = fmt.Sprintf("{id:%s}", "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}")

// MockServer for CLI unit testing
func MockServer(t *testing.T) *httptest.Server {
	policyIdExpr := fmt.Sprintf("%s%s", "/management/v1/policies/", idReg)
	tenantUserExpr := fmt.Sprintf("%s", "/management/v1/users")
	tenantUserIdExpr := fmt.Sprintf("%s%s", "/management/v1/users/", idReg)

	serviceExpr := fmt.Sprintf("%s", "/management/v1/services")
	serviceIdExpr := fmt.Sprintf("%s%s", "/management/v1/services/", idReg)

	subscriptionExpr := fmt.Sprintf("%s%s%s", "/management/v1/services/", idReg, "/api-clients")
	subscriptionIdExpr := fmt.Sprintf("%s%s%s%s", "/management/v1/services/", idReg, "/api-clients/", idReg)
	subscriptionPolicyExpr := fmt.Sprintf("%s%s%s%s%s", "/management/v1/services/", idReg, "/api-clients/", idReg, "/policies")
	subscriptionTagExpr := fmt.Sprintf("%s%s%s%s%s", "/management/v1/services/", idReg, "/api-clients/", idReg, "/tags")

	serviceOfferExpr := fmt.Sprintf("/management/v1/service-offers")

	productExpr := fmt.Sprintf("%s%s%s", "/management/v1/service-offers/", idReg, "/products")

	tenantTagsExpr := fmt.Sprintf("%s", "/management/v1/tags")

	r := mux.NewRouter()

	r.HandleFunc("/management/v1/policies", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(policy))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPost)

	r.HandleFunc("/management/v1/policies", func(w http.ResponseWriter, r *http.Request) {
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

	r.HandleFunc(tenantUserExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(user))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPost)

	r.HandleFunc(tenantUserExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(userList))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(tenantUserIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(user))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(tenantUserIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write(nil)
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodDelete)

	r.HandleFunc(tenantUserIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(user))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPut)

	r.HandleFunc(subscriptionExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(subscription))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPost)

	r.HandleFunc(subscriptionIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(subscription))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPut)

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

	r.HandleFunc(subscriptionIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(subscription))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodDelete)

	r.HandleFunc(subscriptionPolicyExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(policyIds))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(subscriptionTagExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(tagsValues))
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

	r.HandleFunc(serviceIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write(nil)
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodDelete)

	r.HandleFunc(serviceIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Accept", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(service))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPut)

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

	r.HandleFunc(tenantTagsExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Accept", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(tag))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPost)

	r.HandleFunc(tenantTagsExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(tagList))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	return httptest.NewServer(r)
}

//SetupMockConfiguration setting up mock CLI configurations
func SetupMockConfiguration(serverUrl string, configFile *os.File) *config.Configuration {

	c := &config.Configuration{
		AmberBaseUrl: serverUrl,
		TenantId:     "f04971b7-fb41-4a9e-a06e-4bf6e71f98b3",
	}

	fileInfo, _ := os.Stat(configFile.Name())
	if fileInfo.Size() == 0 {
		err := yaml.NewEncoder(configFile).Encode(c)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return c
}

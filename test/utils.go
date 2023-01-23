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
    "policy_name": "Sample_Policy_SGX",
    "policy_type": "Appraisal policy",
    "attestation_type": "SGX Attestation",
    "service_offer_id": "e8a72b7e-c4b1-4bdc-bf40-68f23c68a2aa" }`

	policyList = `[{
    "policy_id": "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
    "policy": "default matches_sgx_policy = false \n\n matches_sgx_policy = true { \n \n input.amber_sgx_isvsvn == 0 \n\n } ",
	"policy_name": "Sample_Policy_SGX",
    "policy_type": "Appraisal policy",
    "attestation_type": "SGX Attestation",
    "service_offer_id": "e8a72b7e-c4b1-4bdc-bf40-68f23c68a2aa" }]`

	user = `{
        "id": "23011406-6f3b-4431-9363-4e1af9af6b13",
        "email": "arijitgh@gmail.com",
        "tenant_roles": [
            {
                "tenant_id": "89120415-6fbc-41c7-b9f2-3b4ba10e87c9",
                "role": {
                        "id": "66ec2e33-8cd3-42b1-8963-c7765205446e",
                        "name": "Tenant Admin"
                    }
            }
        ],
        "active": false,
        "created_at": "2022-06-19T20:02:55.157679Z"
    }`

	userList = `[
    {
        "id": "23011406-6f3b-4431-9363-4e1af9af6b13",
        "email": "arijitgh@gmail.com",
        "role": {
				"id": "66ec2e33-8cd3-42b1-8963-c7765205446e",
				"name": "Tenant Admin"
        },
        "active": false,
        "created_at": "2022-06-19T20:02:55.157679Z"
    },
    {
        "id": "43011406-6f3b-4431-9363-4e1af9af6b11",
        "email": "anotheremail@gmail.com",
        "role": {
				"id": "66ec2e33-8cd3-42b1-8963-c7765205446e",
				"name": "Tenant Admin"
        },
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

	apiClientList = `[
    {
        "id": "3780cc39-cce2-4ec2-a47f-03e55b12e259",
        "service_id": "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
        "product_id": "e169d34f-58ce-4717-9b3a-5c66abd33417",
        "status": "",
        "description": "Test apiClient"
    }]`

	apiClient = `{
        "id": "3780cc39-cce2-4ec2-a47f-03e55b12e259",
        "service_id": "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
        "product_id": "e169d34f-58ce-4717-9b3a-5c66abd33417",
        "status": "",
        "name": "Test apiClient",
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

	plans = `[
		{
			"id": "8f2a20fa-b08d-48a8-b2b4-2ebd1feb6f74",
    		"service_offer_id": "ee28f3c2-6f58-489d-aa46-1140565d4718",
    		"name": "Basic",
			"max_key": 1,
    		"max_tenant_admin": 1,
    		"max_tenant_user": 1,
    		"max_policy": 1,
    		"ledger": false
  		}
	]`

	planProducts = `{
			"id": "8f2a20fa-b08d-48a8-b2b4-2ebd1feb6f74",
  			"service_offer_id": "ee28f3c2-6f58-489d-aa46-1140565d4718",
  			"name": "Premium",
  			"max_key": 10,
  			"max_tenant_admin": 10,
  			"max_tenant_user": 10,
  			"max_policy": 10,
  			"ledger": false,
  			"products": [
			{
      			"id": "a3ad72aa-86d6-49aa-b851-6a52caf0941b",
      			"service_offer_id": "ee28f3c2-6f58-489d-aa46-1140565d4718",
      			"name": "Premium",
      			"policy": {
        			"limit": 40,
        			"quota": 2500000,
        			"limit_renewal_period": 60,
        			"quota_renewal_period": 2592000
      		},
			"plan_id": "838c41b7-8c75-466e-9753-d6c3424662f2",
      		"product_type": ""
		}
 	  ]
	}`
)

var idReg = fmt.Sprintf("{id:%s}", "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}")

// MockServer for CLI unit testing
func MockServer(t *testing.T) *httptest.Server {
	policyIdExpr := fmt.Sprintf("%s%s", "/management/v1/policies/", idReg)
	tenantUserExpr := fmt.Sprintf("%s", "/management/v1/users")
	tenantUserIdExpr := fmt.Sprintf("%s%s", "/management/v1/users/", idReg)

	serviceExpr := fmt.Sprintf("%s", "/management/v1/services")
	serviceIdExpr := fmt.Sprintf("%s%s", "/management/v1/services/", idReg)

	apiClientExpr := fmt.Sprintf("%s%s%s", "/management/v1/services/", idReg, "/api-clients")
	apiClientIdExpr := fmt.Sprintf("%s%s%s%s", "/management/v1/services/", idReg, "/api-clients/", idReg)
	apiClientPolicyExpr := fmt.Sprintf("%s%s%s%s%s", "/management/v1/services/", idReg, "/api-clients/", idReg, "/policies")
	apiClientTagExpr := fmt.Sprintf("%s%s%s%s%s", "/management/v1/services/", idReg, "/api-clients/", idReg, "/tags")

	serviceOfferExpr := fmt.Sprintf("/management/v1/service-offers")

	productExpr := fmt.Sprintf("%s%s%s", "/management/v1/service-offers/", idReg, "/products")

	tenantTagsExpr := fmt.Sprintf("%s", "/management/v1/tags")

	planExpr := fmt.Sprintf("%s%s%s", "/management/v1/service-offers/", idReg, "/plans")
	planIdExpr := fmt.Sprintf("%s%s%s%s", "/management/v1/service-offers/", idReg, "/plans/", idReg)

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

	r.HandleFunc(apiClientExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(apiClient))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPost)

	r.HandleFunc(apiClientIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(apiClient))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPut)

	r.HandleFunc(apiClientExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(apiClientList))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(apiClientIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(apiClient))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(apiClientIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(apiClient))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodDelete)

	r.HandleFunc(apiClientPolicyExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(policyIds))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(apiClientTagExpr, func(w http.ResponseWriter, r *http.Request) {
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

	r.HandleFunc(planExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(plans))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(planIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(planProducts))
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

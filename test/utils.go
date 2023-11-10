/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"intel/tac/v1/config"
	"intel/tac/v1/models"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	policy = `{
        "policy_id": "52135615-3881-4b94-91ff-49f01e626e7b",
        "policy": "default matches_sgx_policy = true \n matches_sgx_policy = true { \n input.sgx_mrenclave == \"83f4e819861adef6ffb2a4865efea9337b91ed30fa33491b17f0d5d9e8204423\" \n input.sgx_mrsigner == \"83d719e77deaca1470f6baf62a4d774303c899db69020f9c70ee1dfc08c7ce9f\" \n }",
        "policy_name": "test-custom4",
        "policy_type": "Appraisal policy",
        "service_offer_id": "736c0281-f13f-4013-a1d2-7ecb92bca0e3",
        "attestation_type": "SGX Attestation",
        "creator_id": "cbb41ee7-e2d2-4b3e-b9f1-37f686f65b34",
        "updater_id": "cbb41ee7-e2d2-4b3e-b9f1-37f686f65b34",
        "deleted": false,
        "created_time": "2023-11-03T12:15:38.181923+05:30",
        "modified_time": "2023-11-03T12:18:11.324105+05:30",
        "policy_hash": "z4yZUmv2fccp8jiRkY3dtkotfDrrCkVSXKRli1O6IKmIS7zlm5qfHJJIGnNpKasP",
        "policy_signature": "hMoTHsLKDRuKkTYOjkv04paXTJmRTEbUDNS9BJcq86iGudgb3j9gUbBGsXH0H7dpz0+5KmuigEbPttan6O/RmwA4xXGrYCQtG4+NB9eM/amNgWNPW7gL482XA18BgJdu4I/f+XGWPbraq+1X2U7yGRi4ise0KM75mTAPdnUBBKq6sk8TgYJGH7n8+ohKbrFyz/IRHvKrwcmnsw7T55ZqU8p1M3N4KSmLz1onKRUXkM5V/6NNYRWFFNytlDFD9P7xU6ns/Ix/FQ1oHCP5WPmTPbnlb512Fy9Ugv2uFr1XUplocjSdL0D18WCedzXcsUwrRrPy2b11drB8TCDMeDVrX+x7zU60r7xICjFaHLJ+K1qDBnfCEQb8uy9XhR/VjQixDZWtAk/hSdCFnOssi+Agrcv3OesXCda1a6yFoinYXDxJqMjryi5vow04yhPkNzR8iiUnIuKLZfTMTM7GRauKXbi2BwA6LWucxa2rhyanuqqO4WB48GJOzIxMYTqz13nP",
        "version": "v2",
        "signed_by_tenant": false
    }`

	policyList = `[{
        "policy_id": "52135615-3881-4b94-91ff-49f01e626e7b",
        "policy": "default matches_sgx_policy = true \n matches_sgx_policy = true { \n input.sgx_mrenclave == \"83f4e819861adef6ffb2a4865efea9337b91ed30fa33491b17f0d5d9e8204423\" \n input.sgx_mrsigner == \"83d719e77deaca1470f6baf62a4d774303c899db69020f9c70ee1dfc08c7ce9f\" \n }",
        "policy_name": "test-custom4",
        "policy_type": "Appraisal policy",
        "service_offer_id": "736c0281-f13f-4013-a1d2-7ecb92bca0e3",
        "attestation_type": "SGX Attestation",
        "creator_id": "cbb41ee7-e2d2-4b3e-b9f1-37f686f65b34",
        "updater_id": "cbb41ee7-e2d2-4b3e-b9f1-37f686f65b34",
        "deleted": false,
        "created_time": "2023-11-03T12:15:38.181923+05:30",
        "modified_time": "2023-11-03T12:18:11.324105+05:30",
        "policy_hash": "z4yZUmv2fccp8jiRkY3dtkotfDrrCkVSXKRli1O6IKmIS7zlm5qfHJJIGnNpKasP",
        "policy_signature": "hMoTHsLKDRuKkTYOjkv04paXTJmRTEbUDNS9BJcq86iGudgb3j9gUbBGsXH0H7dpz0+5KmuigEbPttan6O/RmwA4xXGrYCQtG4+NB9eM/amNgWNPW7gL482XA18BgJdu4I/f+XGWPbraq+1X2U7yGRi4ise0KM75mTAPdnUBBKq6sk8TgYJGH7n8+ohKbrFyz/IRHvKrwcmnsw7T55ZqU8p1M3N4KSmLz1onKRUXkM5V/6NNYRWFFNytlDFD9P7xU6ns/Ix/FQ1oHCP5WPmTPbnlb512Fy9Ugv2uFr1XUplocjSdL0D18WCedzXcsUwrRrPy2b11drB8TCDMeDVrX+x7zU60r7xICjFaHLJ+K1qDBnfCEQb8uy9XhR/VjQixDZWtAk/hSdCFnOssi+Agrcv3OesXCda1a6yFoinYXDxJqMjryi5vow04yhPkNzR8iiUnIuKLZfTMTM7GRauKXbi2BwA6LWucxa2rhyanuqqO4WB48GJOzIxMYTqz13nP",
        "version": "v2",
        "signed_by_tenant": false
    }]`

	user = `{
        "id": "23011406-6f3b-4431-9363-4e1af9af6b13",
        "email": "arijitgh@gmail.com",
		"role": {
			"id": "66ec2e33-8cd3-42b1-8963-c7765205446e",
			"name": "Tenant Admin"
		},
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
        "name": "Test Service"
    }]`
	service = `{
        "id": "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
        "service_offer_id": "ae3d7720-08ab-421c-b8d4-1725c358f03e",
		"service_offer_name": "Test Service Offer",
		"plan_id": "bc3d7720-08ab-421c-b8d4-1725c358f03e",
		"plan_name": "Test Plan",
        "name": "Test Service"
    }`

	apiClientList = `[
    {
        "id": "3780cc39-cce2-4ec2-a47f-03e55b12e259",
        "service_id": "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
        "product_id": "e169d34f-58ce-4717-9b3a-5c66abd33417",
        "status": "",
        "name": "Test apiClient"
    }]`

	apiClientDetails = `{
        "id": "3780cc39-cce2-4ec2-a47f-03e55b12e259",
        "service_id": "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
		"service_offer_name": "Service Offer Name",
        "product_id": "e169d34f-58ce-4717-9b3a-5c66abd33417",
        "status": "",
        "name": "Test apiClient",
		"keys": [
			"9dca50986c414304a4b1ffe202dcf2b0",
			"996a9a6e67814f1784eadb5405bdabf3"
		],
		"policy_ids": [
          	"1cf3db7d-81ea-4904-babf-fcb3501492db",
          	"2cf3db7d-81ea-4904-babf-fcb3501492db",
          	"3cf3db7d-81ea-4904-babf-fcb3501492db",
          	"4cf3db7d-81ea-4904-babf-fcb3501492db"
    	],
    	"tags": [
       	{
            "key": "Workload",
            "value": "Workload-Binary",
            "predefined": true
        }
    	],
     	"created_at": "0001-01-01T00:00:00Z"
    }`

	apiClient = `{
        "id": "3780cc39-cce2-4ec2-a47f-03e55b12e259",
        "service_id": "5cfb6af4-59ac-4a14-8b83-bd65b1e11777",
        "product_id": "e169d34f-58ce-4717-9b3a-5c66abd33417",
        "status": "",
        "name": "Test apiClient"
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
	    "tags": [
         {
            "key": "Workload",
            "value": "AI",
			"predefined": true
		 },
         {
            "key": "Workload",
            "value": "EXE Workload",
			"predefined": false
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

	tenantSettings = `{
			"attest_failure_email" : "dummy@email.com"
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
	tenantTagIdExpr := fmt.Sprintf("%s%s%s", tenantTagsExpr, "/", idReg)

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
		_, err := w.Write([]byte(apiClientDetails))
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
		_, err := w.Write([]byte(apiClientDetails))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(apiClientIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write(nil)
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

	r.HandleFunc(tenantTagIdExpr, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(tagList))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodDelete)

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

	r.HandleFunc("/management/v1/tenants/settings", func(w http.ResponseWriter, r *http.Request) {
		request, _ := io.ReadAll(r.Body)
		dec := json.NewDecoder(bytes.NewReader(request))
		dec.DisallowUnknownFields()
		var settings models.AttestationFailureEmail
		err := dec.Decode(&settings)
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to decode data")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err = w.Write(request)
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodPut)

	r.HandleFunc("/management/v1/tenants/settings", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		_, err := w.Write([]byte(tenantSettings))
		if err != nil {
			t.Log("test/test_utility:mockServer(): Unable to write data")
		}
	}).Methods(http.MethodGet)

	return httptest.NewServer(r)
}

// SetupMockConfiguration setting up mock CLI configurations
func SetupMockConfiguration(serverUrl string, configFile *os.File) *config.Configuration {

	c := &config.Configuration{
		TrustAuthorityBaseUrl: serverUrl,
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

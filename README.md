# Amber CLI
A CLI tool for Tenants to use and access Amber Services

### OS Supported
Ubuntu LTS 20.04

### Build

- git clone https://github.com/intel/amber-cli.git cli
- cd cli and run "make installer"
- copy the binary installer tenantctl-{version}.bin to the system where it needs to be deployed
- create an env file tac.env and add the following mandatory contents:<br>
  a. AMBER_BASE_URL=< URL of API Gateway > <br>
  b. TENANT_ID="< Id of the Tenant >"  (The Tenant Id can be overridden from CLI) <br>
- run "./tenantctl-{version}.bin". This will install the CLI to your system.
- use the CLI: tenantctl < command > < resource >

Note: If behind a proxy, add the Amber FQDN to NO_PROXY environment variable.

### Uninstall 
- run "tenantctl uninstall"

### Setup configuration
- tenantctl config -v < env file path >

### Bash Completion
- tenantctl completion

### Version
- tenantctl version

### Commands Usage examples (please see help for more details ):

##### Create User:
tenantctl create user -a < api key > -e < email Id> -r < Role (Tenant Admin/User) >

##### Get Users:               
tenantctl list user -a < api key >

##### Update User Role:
tenantctl update user role -a < api key > -u < user id > -r < Role (Tenant Admin/User) >

##### Delete User:
tenantctl delete user -a < api key > -u < user id >

##### Get Service Offers:
tenantctl list serviceOffer -a < api key >

##### Get Products:            
tenantctl list product -a < api key > -r < service offer id >

##### Create Service:          
tenantctl create service -a < api key > -r < service offer id > -n < service name >

##### Get Services:
tenantctl list service -a < api key >

##### Get Service By Id:
tenantctl list service -a < api key > -r < service Id >

##### Delete Service:
tenantctl delete service -a < api key > -s < service id >

##### Update Service:
tenantctl update service -a < api key > -s < service id > -n < service name >

##### Create Subscription:
tenantctl create subscription -a < api key > -r < service id > -p < product id > -d < subscription name > -i "comma separated policy Ids" -v "tag-id1:tag-name1,tag-id2:tag-name2"

##### Update Subscription
tenantctl update subscription -a < api key > -r < service id > -p < product id > -u < subscription id > -d < subscription name > -i "comma separated policy Ids" -v "tag-id1:tag-name1,tag-id2:tag-name2" -s < active/inactive >

##### Get Subscriptions:
tenantctl list subscription -a < api key > -r < service id >

##### Get Subscription by id:
tenantctl list subscription -a < api key > -r < service id > -d < subscription id >

##### Delete a subscription:
tenantctl delete subscription -a < api key > -r < service id > -d < subscription id >

##### Create tag:
tenantctl create tag -a < api key > -n < tag name > -t < tenant Id >

##### List tags:
tenantctl list tag -a < api key >

##### List Subscription Policies:
tenantctl list subscription policy -a < api key > -r < service offer id > -s < subscription id >

##### List Subscription Tags:
tenantctl list subscription tag -a < api key > -r < service offer id > -s < subscription id >

##### Create Policy:
tenantctl create policy -a < api key > -f < policy file path >

-  Sample policy for policy create command:

```json
{
  "policy": "default matches_sgx_policy = false \n\n matches_sgx_policy = true { \n input.amber_sgx_is_debuggable == false \n input.amber_sgx_isvsvn == 0 \n input.amber_sgx_isvprodid == 0 \n input.amber_sgx_mrsigner ==  \"d412a4f07ef83892a5915fb2ab584be31e186e5a4f95ab5f6950fd4eb8694d7b\" \n  \n input.amber_sgx_mrenclave == \"bab91f200038076ac25f87de0ca67472443c2ebe17ed9ba95314e609038f51ab\" \n }",
  "user_id": "f04971b7-fb41-4a9e-a06e-4bf6e71f98b3",
  "policy_name": "Sample_Policy_SGX",
  "policy_type": "Appraisal policy",
  "service_offer_name": "SGX Attestation",
  "service_offer_id": "b04971b7-fb41-4a9e-a06e-4bf6e71f98bd"
}
```

##### Get policies:
tenantctl list policies -a < api key >

##### Get policy by id:
tenantctl list policies -a < api key > -p < policy id >

##### Delete policy:
tenantctl delete policy -a < api key > -p < policy id >

##### Update policy:
tenantctl update policy -a < api key > -f < policy file path >

- Sample policy for policy update command:

```json
{
  "policy_id": "e48dabc5-9608-4ff3-aaed-f25909ab9de1",
  "policy": "default matches_sgx_policy = false \n\n matches_sgx_policy = true { \n input.amber_sgx_is_debuggable == false \n input.amber_sgx_isvsvn == 0 \n input.amber_sgx_isvprodid == 0 \n input.amber_sgx_mrsigner ==  \"d412a4f07ef83892a5915fb2ab584be31e186e5a4f95ab5f6950fd4eb8694d7b\" \n  \n input.amber_sgx_mrenclave == \"bab91f200038076ac25f87de0ca67472443c2ebe17ed9ba95314e609038f51ab\" \n }",
  "user_id": "f04971b7-fb41-4a9e-a06e-4bf6e71f98b3",
  "policy_name": "Sample_Policy_SGX",
  "policy_type": "Appraisal policy",
  "service_offer_name": "SGX Attestation",
  "service_offer_id": "b04971b7-fb41-4a9e-a06e-4bf6e71f98bd"
}
```

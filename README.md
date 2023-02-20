# Amber CLI
A CLI tool for Tenants to use and access Amber Services

### OS Supported
Ubuntu LTS 20.04

### Prerequisites
1. Installing make and makeself
   - Run the following command:
   ```console
   apt -y install make makeself
   ```
2. Installing golang
   - Supported golang version is 1.18.x
   - Go through the following link to install golang: https://go.dev/doc/install

### Build

- git clone https://github.com/intel/amber-cli.git cli
- cd cli and run "make installer"
- copy the binary installer tenantctl-{version}.bin to the system where it needs to be deployed
- create an env file tac.env and add the following mandatory contents:<br>
  a. AMBER_BASE_URL=< URL of API Gateway > <br>
  b. TENANT_ID="< Id of the Tenant >"  (The Tenant Id can be overridden from CLI) <br>
- run "./tenantctl-{version}.bin". This will install the CLI to your system.
- use the CLI: tenantctl < command > < resource >
- export AMBER_API_KEY=<subscription key>. This needs to be done before running any of the CLI commands.

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
tenantctl create user -e < email Id> -r < Role (Tenant Admin/User) >

##### Get Users:               
tenantctl list user

##### Get Users by email id:
tenantctl list user -e <email id>

##### Update User Role:
tenantctl update user role -u < user id > -r < Role (Tenant Admin/User) >

##### Delete User:
tenantctl delete user -u < user id >

##### Get Service Offers:
tenantctl list serviceOffer

##### Get Plans:
tenantctl list plan -r < service offer id >

##### Get Plan By Id:
tenantctl list plan -r < service offer id > -p < plan id >

##### Get Products:            
tenantctl list product -r < service offer id >

##### Get Services:
tenantctl list service

##### Get Service By Id:
tenantctl list service -r < service Id >

##### Create Api Client:
tenantctl create apiClient -r < service id > -p < product id > -n < api client name > -i "comma separated policy Ids" -v "tag-key1:tag-value1,tag-key2:tag-value2"

##### Update Api Client:
tenantctl update apiClient -r < service id > -p < product id > -c < api client id > -n < api client name > -i "comma separated policy Ids" -v "tag-key1:tag-value1,tag-key2:tag-value2" -s < Active/Inactive/Cancelled >

##### Get Api Clients:
tenantctl list apiClient -r < service id >

##### Get Api Client by id:
tenantctl list apiClient -r < service id > -c < api client id >

##### Delete an Api Client:
tenantctl delete apiClient -r < service id > -c < api client id >

##### Create tag:
tenantctl create tag -n < tag name > -t < tenant Id >

##### List tags:
tenantctl list tag

##### List Api Client Policies:
tenantctl list apiClient policy -r < service id > -c < api client id >

##### List Api Client Tags:
tenantctl list apiClient tag -r < service id > -c < api client id >

##### Create Policy JWT
tenantctl create policy-jwt -f <rego policy file path> -p <signing key path> -c <cert path> -a <algorithm> -s

##### Create Policy:
tenantctl create policy -n < name of policy > -t < policy type > -a < attestation type > -r < service offer id > -f < rego policy file path >

##### Get policies:
tenantctl list policy

##### Get policy by id:
tenantctl list policy -p < policy id >

##### Delete policy:
tenantctl delete policy -p < policy id >

##### Update policy:
tenantctl update policy -i < policy id > -n < name of policy > -f < rego policy file path >

-  Sample rego policy for create/update policy command:

```bash
default matches_sgx_policy = false 
matches_sgx_policy = true 
{  input.amber_sgx_is_debuggable == false 
   input.amber_sgx_isvsvn == 0 
   input.amber_sgx_isvprodid == 0 
   input.amber_sgx_mrsigner ==  \"d412a4f07ef83892a5915fb2ab584be31e186e5a4f95ab5f6950fd4eb8694d7b\" 
   input.amber_sgx_mrenclave == \"bab91f200038076ac25f87de0ca67472443c2ebe17ed9ba95314e609038f51ab\" 
} 
```


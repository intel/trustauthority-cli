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
   - Supported golang version is 1.20.x
   - Go through the following link to install golang: https://go.dev/doc/install

3. Add the local binary path, namely $HOME/.local/bin/, to your PATH environment variable if not already present

### Build

- git clone https://github.com/intel/amber-cli.git cli
- cd cli and run "make installer"
- copy the binary installer tenantctl-{version}.bin to the system where it needs to be deployed
- create an env file tac.env in your home directory and add the following mandatory contents:<br>
  a. AMBER_BASE_URL=< URL of API Gateway > <br>
  b. AMBER_API_KEY="< API Key of the Tenant >" <br>
- run "./tenantctl-{version}.bin". This will install the CLI to your system.
- use the CLI: tenantctl < command > < resource >

Note: If behind a proxy, add the Amber FQDN to NO_PROXY environment variable.

### Directory structure

All files are stored in user's home directory. Following are the details:

- Configuration: $HOME/.config/tenantctl/config.yaml
- Logs: $HOME/.config/tenantctl/logs/tac.log
- Bin: $HOME/.local/bin/tenantctl

Note: If you cannot access the command, add the binary path to the PATH env variable

## Commands

Note: Request ID could be a randomly generated string of at most 128 bytes which can work as a unique 
identifier for each CRUD operation. This can be provided as an optional parameter to all the CRUD commands only.

### Uninstall 
- tenantctl uninstall

### Setup configuration
- tenantctl config -v < env file path >

### Bash Completion
- tenantctl completion

### Version
- tenantctl version

### Commands Usage examples (please see help for more details ):

##### Create User:
tenantctl create user -q < request id > -e < email Id> -r < Role (Tenant Admin/User) >

##### Get Users:               
tenantctl list user -q < request id >

##### Get Users by email id:
tenantctl list user -q < request id > -e <email id>

##### Update User Role:
tenantctl update user role -q < request id > -u < user id > -r < Role (Tenant Admin/User) >

##### Delete User:
tenantctl delete user -q < request id > -u < user id >

##### Delete Tag:
tenantctl delete tag -q < request id > -t < tag id >

##### Get Service Offers:
tenantctl list serviceOffer

##### Get Plans:
tenantctl list plan -q < request id > -r < service offer id >

##### Get Plan By Id:
tenantctl list plan -q < request id > -r < service offer id > -p < plan id >

##### Get Products:            
tenantctl list product -q < request id > -r < service offer id >

##### Get Services:
tenantctl list service -q < request id >

##### Get Service By Id:
tenantctl list service -q < request id > -r < service Id >

##### Create Api Client:
tenantctl create apiClient -q < request id > -r < service id > -p < product id > -n < api client name > -i "comma separated policy Ids" -v "tag-key1:tag-value1,tag-key2:tag-value2"

##### Update Api Client:
tenantctl update apiClient -q < request id > -r < service id > -p < product id > -c < api client id > -i "comma separated policy Ids" -v "tag-key1:tag-value1,tag-key2:tag-value2" -s < Active/Inactive/Cancelled >

##### Get Api Clients:
tenantctl list apiClient -q < request id > -r < service id >

##### Get Api Client by id:
tenantctl list apiClient -q < request id > -r < service id > -c < api client id >

##### Delete an Api Client:
tenantctl delete apiClient -q < request id > -r < service id > -c < api client id >

##### Create tag:
tenantctl create tag -q < request id > -n < tag name >

##### List tags:
tenantctl list tag -q < request id >

##### List Api Client Policies:
tenantctl list apiClient policy -q < request id > -r < service id > -c < api client id >

##### List Api Client Tags:
tenantctl list apiClient tag -q < request id > -r < service id > -c < api client id >

##### Update Tenant Settings:
tenantctl update tenant-settings -q < request id > -e < email id >

##### Update Tenant Settings (disable notification):
tenantctl update tenant-settings -q < request id > -d

##### List Tenant Settings:
tenantctl list tenant-settings -q < request id >

##### Create Policy:
tenantctl create policy -q < request id > -n < name of policy > -t < policy type > -a < attestation type > -r < service offer id > -f < rego policy file path >
Note: Policy file size should be <= 10KB

##### Get policies:
tenantctl list policy -q < request id >

##### Get policy by id:
tenantctl list policy -q < request id > -p < policy id >

##### Delete policy:
tenantctl delete policy -q < request id > -p < policy id >

##### Update policy:
tenantctl update policy -q < request id > -i < policy id > -n < name of policy > -f < rego policy file path >
Note: Policy file size should be <= 10KB

-  Sample rego policy for create/update policy command:

```bash
default matches_sgx_policy = false 
matches_sgx_policy = true 
{  input.amber_tee_is_debuggable == false 
   input.amber_sgx_isvsvn == 0 
   input.amber_sgx_isvprodid == 0 
   input.amber_sgx_mrsigner ==  \"d412a4f07ef83892a5915fb2ab584be31e186e5a4f95ab5f6950fd4eb8694d7b\" 
   input.amber_sgx_mrenclave == \"bab91f200038076ac25f87de0ca67472443c2ebe17ed9ba95314e609038f51ab\" 
} 
```

### Create Policy JWT
tenantctl create policy-jwt -q < request id > -f < rego policy file path > -p < signing key path > -c < cert path > -a < algorithm > -s

#### Prerequisites: 
Create self signed key and certificate for policy JWT token creation:
- Generate key and cert files for -algorithm (PS384 | RS384) (Recommend)
```
openssl req -x509 -nodes -days 365 -newkey rsa:3072 -keyout amber-jwt.key -out amber-jwt.crt
```
- Generate key and cert files for -algorithm (PS256 | RS256)
```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout amber-jwt.key -out amber-jwt.crt
```

#### Notes:
1. Signed policy token could be self verified at jwt.io
2. Output file name of this command is input policy file name suffixed with ".signed.current_timestamp.txt" extension
3. Policy payload Amber uses rego format which is different from Azure MAA
4. Supported signing algorithms are "RS256", "PS256", "RS384", "PS384", default algorithm is PS384
5. The signing algorithm needs to match the certificate algorithm


#### References:
- Azure MAA:
    - https://learn.microsoft.com/en-us/azure/attestation/policy-examples
    - https://docs.microsoft.com/en-us/azure/attestation/author-sign-policy
- JWS RFC7515:
    - https://www.rfc-editor.org/rfc/rfc7515

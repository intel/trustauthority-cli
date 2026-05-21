# Intel Trust Authority CLI
The Intel® Trust Authority CLI is an open-source tool tenants use to make API calls to Intel's Trust Authority. The source code for the Intel® Trust Authority CLI is available on GitHub.

### OS Supported
Ubuntu LTS 20.04

### Prerequisites
The make, makeself, and golang packages are required before installing the Intel Trust Authority CLI. The steps below explain how to install these packages.

1. Install make and makeself.
    - Run the following command to install the make and makeself packages.
   ```console
   apt -y install make makeself
   ```
2. Installing golang.
    - Supported golang version is 1.26.3
    - To install golang, follow the instructions at the following link:https://go.dev/doc/install

3. Add the local binary path, `$HOME/.local/bin/`, to your PATH environment variable if not already present.

### Build the Intel Trust Authority CLI
Follow the steps below to build the Intel Trust Authority CLI.
1. Create a cli directory.
   `mkdir cli`
2. Clone the trustauthority-cli repository by running the following command.
   `git clone https://github.com/intel/trustauthority-cli.git trustauthority-cli`
3. Create the CLI installer in the cli directory.
   `cd trustauthority-cli and run "make installer"`
4. Copy the binary installer trustauthorityctl-{version}.bin to the system where it needs to be deployed.
5. Create an env file trustauthorityctl.env in your home directory and add the following mandatory contents:<br>
   a. TRUSTAUTHORITY_URL="https://api.trustauthority.intel.com" <br>
   b. TRUSTAUTHORITY_API_KEY="< Management API Key of the Tenant >" <br>
6. To install the CLI to your system, run the following command:
   `./trustauthorityctl-{version}.bin`
7. To use the CLI, follow this syntax:
   `trustauthorityctl < command > < resource >`

> [!NOTE]
> If you are in the European Union (EU) region, use the following Intel Trust Authority API URL:
>
> TRUSTAUTHORITY_URL="https://api.eu.trustauthority.intel.com" 

### Directory structure

All files are stored in the user's home directory. The contents of the directories are listed below:

- Configuration: $HOME/.config/trustauthorityctl/config.yaml
- Logs: $HOME/.config/trustauthorityctl/logs/trustauthorityctl.log
- Bin: $HOME/.local/bin/trustauthorityctl

> [!Note]
> If you cannot access the command, add the binary path to the PATH env variable.

## Commands

> [!Note]
> The Request ID could be a randomly generated string of 128 bytes or less, which can be a unique
identifier for each CRUD operation. The Request ID can only be provided as an optional parameter for CRUD commands.

### Uninstall
To uninstall the Intel Trust Authority CLI, run the following command:
`trustauthorityctl uninstall`

### Setup configuration
The file path to the trustauthorityctl.env file created in the previous step is needed to complete this step.

To configure the Intel Trust Authority CLI, run the command below.
`trustauthorityctl config -v < env file path >`

### Bash Completion
To install bash completion for the Intel Trust Authority CLI, run the following command:
`trustauthorityctl completion`

### Version
To get the version number of the tenant CLI installed on your system, run the following command:
`trustauthorityctl version`

### Commands Usage examples (please see help for more details ):

##### Create User:
trustauthorityctl create user -q < request id > -e < email Id> -r < Role (Tenant Admin/User) >

##### Get Users:
trustauthorityctl list user -q < request id >

##### Get Users by email ID:
trustauthorityctl list user -q < request id > -e <email id>

##### Update User Role:
trustauthorityctl update user role -q < request id > -u < user id > -r < Role (Tenant Admin/User) >

##### Delete User:
trustauthorityctl delete user -q < request id > -u < user id >

##### Delete Tag:
trustauthorityctl delete tag -q < request id > -t < tag id >

##### Get Service Offers:
trustauthorityctl list serviceOffer

##### Get Plans:
trustauthorityctl list plan -q < request id > -r < service offer id >

##### Get Plan By Id:
trustauthorityctl list plan -q < request id > -r < service offer id > -p < plan id >

##### Get Products:
trustauthorityctl list product -q < request id > -r < service offer id >

##### Get Services:
trustauthorityctl list service -q < request id >

##### Get Service By Id:
trustauthorityctl list service -q < request id > -r < service Id >

##### Create Api Client:
trustauthorityctl create apiClient -q < request id > -r < service id > -p < product id > -n < api client name > -i "comma separated policy Ids" -v "tag-key1:tag-value1,tag-key2:tag-value2"

##### Update Api Client:
trustauthorityctl update apiClient -q < request id > -r < service id > -p < product id > -c < api client id > -i "comma separated policy Ids" -v "tag-key1:tag-value1,tag-key2:tag-value2" -s < Active/Inactive/Cancelled >

##### Get Api Clients:
trustauthorityctl list apiClient -q < request id > -r < service id >

##### Get Api Client by id:
trustauthorityctl list apiClient -q < request id > -r < service id > -c < api client id >

##### Delete an Api Client:
trustauthorityctl delete apiClient -q < request id > -r < service id > -c < api client id >

##### Create tag:
trustauthorityctl create tag -q < request id > -n < tag name >

##### List tags:
trustauthorityctl list tag -q < request id >

##### List Api Client Policies:
trustauthorityctl list apiClient policy -q < request id > -r < service id > -c < api client id >

##### List Api Client Tags:
trustauthorityctl list apiClient tag -q < request id > -r < service id > -c < api client id >

##### Update Tenant Settings:
trustauthorityctl update tenant-settings -q < request id > -e < email id >

##### Update Tenant Settings (disable notification):
trustauthorityctl update tenant-settings -q < request id > -d

##### List Tenant Settings:
trustauthorityctl list tenant-settings -q < request id >

##### Create a RIM:
trustauthorityctl create rim -q < request id > -n < RIM name > -f < RIM content file path >

Examples:
- Create RIM with plain JSON content:
  `trustauthorityctl create rim -n "acme.rims.mrtd" -f ./rim-content.json`
- Create RIM with signed JWT content:
  `trustauthorityctl create rim -n "public.acme.rim.signed" -f ./rim-signed.txt`

Note: 
- RIM file size should be <= 20KB
- RIM name format: Use dot-separated namespaced names (e.g., "acme.rims.mrtd")
- Public RIMs: Prefix with "public." (e.g., "public.acme.rims.mrtd") to make them accessible across tenants
- Content can be either plain JSON or signed JWT format

##### List all RIMs:
trustauthorityctl list rim -q < request id >

##### Get specific RIM by ID:
trustauthorityctl list rim -q < request id > -r < RIM id >

##### Search RIMs by name:
trustauthorityctl list rim -q < request id > -n < RIM name >

##### Update a RIM:
trustauthorityctl update rim -q < request id > -r < RIM id > -f < RIM content file path > -d < description >

Examples:
- Update RIM with plain JSON content:
  `trustauthorityctl update rim -r "1722b479-cbba-4f10-9ac8-c21dff5bf8db" -f ./new-rim-content.json`
- Update RIM with signed JWT content:
  `trustauthorityctl update rim -r "1722b479-cbba-4f10-9ac8-c21dff5bf8db" -f ./rim-signed.txt`
- Update RIM description only:
  `trustauthorityctl update rim -r "1722b479-cbba-4f10-9ac8-c21dff5bf8db" -d "Updated description"`
- Update both content and description:
  `trustauthorityctl update rim -r "1722b479-cbba-4f10-9ac8-c21dff5bf8db" -f ./new-rim-content.json -d "Updated"`

Note: 
- You can update content, description, or both. At least one field must be provided.
- Content can be either plain JSON or signed JWT format (same as create rim)
- File size should be <= 20KB

##### Delete a RIM:
trustauthorityctl delete rim -q < request id > -r < RIM id >

Example:
`trustauthorityctl delete rim -r "1722b479-cbba-4f10-9ac8-c21dff5bf8db"`

##### Create Policy:
trustauthorityctl create policy -q < request id > -n < name of policy > -t < policy type > -a < attestation type > -r < service offer id > -f < rego policy file path >
Note: Policy file size should be <= 20KB

##### Get policies:
trustauthorityctl list policy -q < request id >

##### Get policy by ID:
trustauthorityctl list policy -q < request id > -p < policy id >

##### Delete policy:
trustauthorityctl delete policy -q < request id > -p < policy id >

##### Update policy:
trustauthorityctl update policy -q < request id > -i < policy id > -n < name of policy > -f < rego policy file path >
Note: Policy file size should be <= 20KB

-  Sample rego policy for create/update policy command:

```bash
default matches_sgx_policy = false
matches_sgx_policy = true
{  input.sgx_is_debuggable == false
   input.sgx_isvsvn == 0
   input.sgx_isvprodid == 0
   input.sgx_mrsigner ==  \"d412a4f07ef83892a5915fb2ab584be31e186e5a4f95ab5f6950fd4eb8694d7b\"
   input.sgx_mrenclave == \"bab91f200038076ac25f87de0ca67472443c2ebe17ed9ba95314e609038f51ab\"
}
```

### Create JWT (Also works for RIM content)

> [!IMPORTANT]
> **Command renamed:** `create jwt` replaces `create policy-jwt`.
> The old name `policy-jwt` remains available as a deprecated alias for backward compatibility but will be removed in a future release.

> **Migrate now:** replace `trustauthorityctl create policy-jwt` with `trustauthorityctl create jwt` in all scripts and workflows.

```
trustauthorityctl create jwt -q < request id > -f < rego policy or RIM content file path > -p < signing key path > -c < cert path > -a < algorithm > -s
```

#### Prerequisites:
Create a self-signed key and certificate for JWT token creation:
- Generate key and cert files for -algorithm (PS384 | RS384) (Recommend)
```
openssl req -x509 -nodes -days 365 -newkey rsa:3072 -keyout ta-jwt.key -out ta-jwt.crt
```
- Generate key and cert files for -algorithm (PS256 | RS256)
```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ta-jwt.key -out ta-jwt.crt
```

#### Usage Examples:
1. Create signed JWT for a rego policy:
   ```bash
   trustauthorityctl create jwt -f ./my-policy.rego -p ./ta-jwt.key -c ./ta-jwt.crt -a PS384 --sign
   ```

2. Create signed JWT for RIM content:
   ```bash
   trustauthorityctl create jwt -f ./rim-content.json -p ./ta-jwt.key -c ./ta-jwt.crt -a PS384 --sign
   ```

3. Create unsigned JWT (algorithm=none):
   ```bash
   trustauthorityctl create jwt -f ./my-policy.rego
   ```

> [!NOTE]
> **Deprecated alias:** The following form still works but is deprecated — update your scripts to use `create jwt`.
> ```bash
> trustauthorityctl create policy-jwt -f ./my-policy.rego -p ./ta-jwt.key -c ./ta-jwt.crt -a PS384 --sign
> ```

#### Notes:
1. The signed policy/RIM token can be self-verified at jwt.io.
2. The output file name of this command is the input file name suffixed with the ".signed.current_timestamp.txt" extension.
3. For policies: The payload uses the rego format, which is different from Azure MAA.
4. For RIMs: The payload should be valid JSON containing RIM measurements and metadata.
5. Supported signing algorithms are "RS256", "PS256", "RS384", "PS384", and the default algorithm is PS384.
6. The signing algorithm needs to match the certificate algorithm.
7. The generated JWT can be used directly when creating or updating policies/RIMs.


#### References:
- Azure MAA:
    - https://learn.microsoft.com/en-us/azure/attestation/policy-examples
    - https://docs.microsoft.com/en-us/azure/attestation/author-sign-policy
- JWS RFC7515:
    - https://www.rfc-editor.org/rfc/rfc7515

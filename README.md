# applications.security.amber.cli
A CLI tool for Tenants to use to access Amber Services

### OS Supported
Ubuntu LTS 20.04

### Build

- git clone https://github.com/intel-innersource/applications.security.amber.core-services.git core-services
- cd core-services/cli and run "make installer"
- copy the binary installer tenantctl-{version}.bin to the system where it needs to be deployed
- create an env file tac.env and add the following mandatory contents:<br>
  a. AMBER_BASE_URL=< URL of API Gateway > <br>
  b. TENANT_ID="< Id of the Tenant >"  (The Tenant Id can be overridden from CLI) <br>
- run "./tenantctl-{version}.bin". This will install the CLI to your system.
- use the CLI: tenantctl < command > < resource >

### Uninstall 
- run "tenantctl uninstall"

### Commands Usage examples (please see help for more details ):

##### Create User:
tenantctl create user -a < api key > -e < email Id> -r < Role (Tenant Admin/User) >

##### Get Users:               
tenantctl list user -a < api key >

##### Get User by id:          
tenantctl list user -a < api key > -u < user id >

##### Update User:
tenantctl update user -a < api key > -u < user id >

##### Update User Role:
tenantctl update user role -a < api key > -u < user id > -r "Tenant Admin,User"

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

##### Delete Service:
tenantctl update service -a < api key > -s < service id >

##### Update Service:
tenantctl update service -a < api key > -s < service id > -n < service name >

##### Create Subscription:
tenantctl create subscription -a < api key > -r < service offer id > -p < product id > -d 'Arijit with tags Subscription' -i 095567af-75e0-4ed3-baa9-74b5242e3061,d40a9a17-8dab-4aed-a465-b7a404ad0ac5 -v "4f3cdd0d-4a2b-445f-9277-07bed3cea808:Workload,42501785-4234-42a3-9ea6-bf3c58244cad:AI"

##### Get Subscriptions:
tenantctl list subscription -a < api key > -r < service offer id >

##### Get Subscription by id:
tenantctl list subscription -a < api key > -r < service offer id > -d < subscription id >

##### Create tag:
tenantctl create tag -a < api key > -n "Arijit tag" -t 5aeb5c92-a6c7-4741-86c1-c8ec8849ed36

##### List tags:
tenantctl list tag -a < api key >

##### List Subscription Policies:
tenantctl list subscription policy -a < api key > -r < service offer id > -s < subscription id >

##### List Subscription Tags:
tenantctl list subscription tag -a < api key > -r < service offer id > -s < subscription id >

##### Create Policy:
tenantctl create policy -a < api key > -f sample/policy.json

##### Get policies:
tenantctl list policies -a < api key >

##### Get policy by id:
tenantctl list policies -a < api key > -p < policy id >

##### Delete policy:
tenantctl delete policy -a < api key > -p < policy id >

##### Update policy:
tenantctl update policy -a < api key > -p < policy id > -f sample/policy.json

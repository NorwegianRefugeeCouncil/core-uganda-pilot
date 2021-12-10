# Core

**Core** is a set of components that attempt to solve the problem of

- **data collection**
- **beneficiary identity**
- **case management**
- **reporting**

for the humanitarian sector.

Humanitarian organizations offer services to beneficiaries, usually in multiple countries. These services might be the

- delivery of goods
- delivery of services, such as education, legal help, medical help, ...

Humanitarian organizations need to collect data about the recipients of those services to

- Offer services to those beneficiaries
- Track the status of the services
- Evaluate the admissibility of a beneficiary to these services
- Accountability to the donors

Also, humanitarian organizations often need to collect other types of data to

- Evaluate where aid efforts should be focused
- Evaluate the status of a population

Humanitarian organizations usually use something like ODK or KOBO to perform this task. While ODK/Kobo and such tools
perform well for offline data collection, The **Norwegian Refugee Council** found some shortcomings with these tools
such as

- Lacks the concept of identity or cases
- Lacks OIDC/OAuth features
- Lack of centralized identity management
- Lack of centralized permission management
- Lack of case management features
- Beneficiaries are not allowed to log in

## Applications

| Name | Description |
|------|-------------|
`core-frontend` | Frontend application that allows managing `forms` and `records`
`login-frontend` | Server Application that allows login and identity provider federation
`forms-api` | API Server for managing `forms` and `records`
`authnz-api` | API Server for managing `OAuth` clients
`authnz-frontend` | Frontend application to manage authorization/authentication
`authnz-bouncer` | Application that verifies, authenticates and authorizes requests

## Dependencies

| Name | Description |
|------|-------------|
`ory/hydra` | OIDC Wrapper for the `login-frontend`. Acts as the Identity Provider for all apps
`spicedb` | Stores authorization entries and evaluates if actions are authorized or not  

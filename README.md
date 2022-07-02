![Gopher Worker](docs/icon.png)

# SupervisorAPI
Task managing API for managers and technicians

## Contents     

- [What is Supervisor API](#what-is-supervisor-api) 
    - [Summary](#summary) 
    - [Features](#features) 
- [Requirements](#requirements) 
- [Instructions](#instructions) 
    - [Auth0 Integration](#auth0-integration) 
    - [Environment variables](#environment-variables) 
    - [Auth2 Client Configuration](#auth2-client-configuration) 
    - [Setting Up Local Environment](#setting-up-local-environment) 
    - [Run project](#run-project) 
- [Available Endpoints](#available-endpoints) 
- [Testing and Coverage](#testing-and-coverage)

# What is Supervisor API
## Summary
Supervisor API is an API that allows CRUD operations on Tasks. This API leverages role based credentials configured and managed in an external Oauth2 authentication service has access permissions.
Existing roles and allowed actions are differentiated by access control:
- Manager
    - Create tasks
    - Update own tasks
    - Fetch any task by task identifier
    - List any tasks by query parameters
- Technician
    - Create tasks
    - Update own tasks
    - Fetch own task by task identifier
    - List own tasks by query parameters

    ## Features
    (needs work)
    [x] Task summary is constraint to 2500 characters
    [x] Task summary is encrypted on database
    [x] Task List endpoints query by "worker_name"
    [x] Task List endpoints query by date "after"
    [x] Task List endpoints query by date "before"
    [ ] Use message broker 

# Requirements
1. Docker Destop running
2. Go version 1.17+
# Instructions

## Auth0 integration
This API uses Auth0, an external authentication service. For testing porpuses, for the default environment variables present in the [.env](https://github.com/MrBolas/SupervisorAPI/blob/b90ce1c7519fe3a813d4515b5aef018027e2f346/.env) the service will use a set of configured test users.

| User          | Password        | Email                   | Role      |
| ------------- |-----------------| ------------------------|-----------|
| Robert        | Robert1234      | robert@supervisor.com   |Manager    |
| Joseph        | Joseph1234      | joseph@supervisor.com   |Technician |
| Cassandra     | Cassandra1234   | cassandra@supervisor.com|Technician |

## Environment variables
Default environment variables are defined in the [.env](https://github.com/MrBolas/SupervisorAPI/blob/b90ce1c7519fe3a813d4515b5aef018027e2f346/.env) file. In this section I'll explain what they are:
| Variable              | Default                                                         | Description                      |
| ----------------------|-----------------------------------------------------------------|----------------------------------|
| MYSQL_USERNAME        | user                                                            | DB username                      |
| MYSQL_PASSWORD        | password                                                        | DB password                      |
| MYSQL_HOSTNAME        | localhost                                                       | DB address                       |
| MYSQL_PORT            | 3306                                                            | DB port                          |
| MYSQL_DATABASE        | sh_supervisor                                                   | DB name                          |
| AUTH0_DOMAIN          | dev-04detuv7.us.auth0.com                                       | auth0 domain address             |
| AUTH0_CLIENT_ID       | fmSt7Lf2b2mQr5LYpsSKglJYMy5YZiJd                                | auth0 clientID                   |
| AUTH0_CLIENT_SECRET   | zw78e03sMF7AqWzQ-ekzTZgqqL93YTxaPwzKxIYNr-KG5aih5eHq2R-rrgy6m-aJ| auth0 client Secret              |
| AUTH0_PUBLIC_KEY_URL  | https://dev-04detuv7.us.auth0.com/.well-known/jwks.json         | auth0 public key                 |
| CRYPTO_KEY            | skidMAhiçWdh34KlosQLP84GhT62smn                                 | Encryption key for sensitive data|
## Auth2 Client Configuration
Any http client used with this service should use a "Resource Owner Password Credentials" Grant type with the following configurations:

| Configuration    | Value (for default Auth0 integration)          |
| -----------------|----------------------------------------        |
| Username         | [User email](#auth0-integration)               |
| Password         | [Password](#auth0-integration)                 |
| Access Token Url | https://dev-04detuv7.us.auth0.com/oauth/token  |
| Client ID        | [AUTH0_CLIENT_ID](#auth0-integration)          |
| Client Secret    | [AUTH0_CLIENT_SECRET](#auth0-integration)      |

## Setting Up Local Environment
1. Launch MySQL docker container (default environment variables already configured)
```
 docker run --name sh_mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password -e MYSQL_PASSWORD=password -e MYSQL_DATABASE=sh_supervisor -e MYSQL_USER=user -d mysql:5.7
```

## Run project
1. Git clone
```
git clone https://github.com/MrBolas/SupervisorAPI.git
```
2. Navigate to folder and install packages
```
go install
```
3. Build
```
go build
```
4. Run it
```
./SupervisorAPI
```

# Available Endpoints
(needs work)

# Testing and Coverage
(needs work)
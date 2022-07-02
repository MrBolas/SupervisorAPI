![Gopher Worker](docs/icon.png)

# SupervisorAPI
Task managing API for managers and technicians

## Contents     

- [What is Supervisor API](#what-is-supervisor-api) 
    - [Summary](#summary) 
    - [Features](#features) 
- [Instructions](#instructions) 
    - [Auth0 Integration](#auth0-integration) 
    - [Auth2 Client Configuration](#auth2-client-configuration) 
    - [Development Environment](#development-environment)
        - [Development Requirements](#development-requirements) 
        - [Project Environment variables](#project-environment-variables) 
        - [Setting Up Local Environment](#setting-up-local-environment) 
        - [Run project](#run-project) 
        - [Monitor Development Redis](#monitor-development-redis)
    - [Docker Compose Environment](#docker-compose-environment)
        - [Docker Compose Requirements](#docker-compose-requirements) 
        - [Docker Environment variables](#docker-environment-variables) 
        - [Setting Up Docker Compose Environment](#setting-up-docker-compose-environment) 
        - [Run Docker Compose](#run-docker-compose) 
        - [Monitor Docker Compose Redis](#monitor-docker-compose-redis)
    - [Kubernetes Environment](#k8s-environment)
        - [K8s Requirements](#k8s-requirements) 
        - [K8s Environment variables](#k8s-environment-variables) 
        - [Setting Up K8s Environment](#setting-up-local-environment) 
        - [Run K8s](#run-k8s) 
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
    [x] Use Redis as message broker 
# Instructions

## Auth0 integration
This API uses Auth0, an external authentication service. For testing porpuses, for the default environment variables present in the [.env](https://github.com/MrBolas/SupervisorAPI/blob/b90ce1c7519fe3a813d4515b5aef018027e2f346/.env) the service will use a set of configured test users.

| User          | Password        | Email                   | Role      |
| ------------- |-----------------| ------------------------|-----------|
| Robert        | Robert1234      | robert@supervisor.com   |Manager    |
| Joseph        | Joseph1234      | joseph@supervisor.com   |Technician |
| Cassandra     | Cassandra1234   | cassandra@supervisor.com|Technician |

## Auth2 Client Configuration
Any http client used with this service should use a "Resource Owner Password Credentials" Grant type with the following configurations:

| Configuration    | Value (for default Auth0 integration)          |
| -----------------|----------------------------------------        |
| Username         | [User email](#auth0-integration)               |
| Password         | [Password](#auth0-integration)                 |
| Access Token Url | https://dev-04detuv7.us.auth0.com/oauth/token  |
| Client ID        | [AUTH0_CLIENT_ID](#auth0-integration)          |
| Client Secret    | [AUTH0_CLIENT_SECRET](#auth0-integration)      |

To falicitate the usage of this API a export of Insomnia endpoints configuration is provided in the [project](https://github.com/MrBolas/SupervisorAPI/blob/6fedd7cbca98cb253a2227fadd771300d221ff37/InsomniaExport/Insomnia_export.json).
## Development Environment
### Development Requirements
1. Docker Destop running
2. Go version 1.17+

### Project Environment variables
Default environment variables are defined in the [.env](https://github.com/MrBolas/SupervisorAPI/blob/6fedd7cbca98cb253a2227fadd771300d221ff37/.env) file. In this section I'll explain what they are:
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
| REDIS_HOST            | localhost                                                       | Redis address                    |
| REDIS_PORT            | 6379                                                            | Redis Port                       |
| REDIS_DB              | 0                                                               | Redis DB                         |

## Setting Up Local Environment
1. Launch MySQL docker container (default environment variables already configured)
```
 docker run --name sh_mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password -e MYSQL_PASSWORD=password -e MYSQL_DATABASE=sh_supervisor -e MYSQL_USER=user -d mysql:5.7
```
2. Launch Redis docker container
```
docker run --name some-redis -d redis
```
## Run project
1. Git clone
```
git clone https://github.com/MrBolas/SupervisorAPI.git
```
2. Navigate to folder and fetch external packages
```
go get -d -v ./...
```
3. Install packages
```
go install -v ./...
```
4. Build
```
go build
```
5. Run executable
```
./SupervisorAPI
```
### Monitor Development Redis
To validate incoming notifications we can use:
```
docker exec -it some-redis redis-cli -h localhost subscribe notifications
```
This command subscribes to the topic notifications in redis-cli inside the docker container some-redis, and will show received messages received for that topic.
## Docker Compose Environment
### Docker Compose Requirements
1. Docker Destop running
### Docker Environment variables
Default environment variables are defined in the [.env-docker](https://github.com/MrBolas/SupervisorAPI/blob/6fedd7cbca98cb253a2227fadd771300d221ff37/.env-docker) file. In this section I'll explain what they are:
| Variable              | Default                                                         | Description                      |
| ----------------------|-----------------------------------------------------------------|----------------------------------|
| MYSQL_USERNAME        | user                                                            | DB username                      |
| MYSQL_PASSWORD        | password                                                        | DB password                      |
| MYSQL_HOSTNAME        | db                                                              | DB address                       |
| MYSQL_PORT            | 3306                                                            | DB port                          |
| MYSQL_DATABASE        | sh_supervisor                                                   | DB name                          |
| AUTH0_DOMAIN          | dev-04detuv7.us.auth0.com                                       | auth0 domain address             |
| AUTH0_CLIENT_ID       | fmSt7Lf2b2mQr5LYpsSKglJYMy5YZiJd                                | auth0 clientID                   |
| AUTH0_CLIENT_SECRET   | zw78e03sMF7AqWzQ-ekzTZgqqL93YTxaPwzKxIYNr-KG5aih5eHq2R-rrgy6m-aJ| auth0 client Secret              |
| AUTH0_PUBLIC_KEY_URL  | https://dev-04detuv7.us.auth0.com/.well-known/jwks.json         | auth0 public key                 |
| CRYPTO_KEY            | skidMAhiçWdh34KlosQLP84GhT62smn                                 | Encryption key for sensitive data|
| REDIS_HOST            | queue                                                           | Redis address                    |
| REDIS_PORT            | 6379                                                            | Redis Port                       |
| REDIS_DB              | 0                                                               | Redis DB                         |

### Setting Up Docker Compose Environment
The deployment orchestrated in the [docker-compose.yaml](https://github.com/MrBolas/SupervisorAPI/blob/6fedd7cbca98cb253a2227fadd771300d221ff37/docker-compose.yaml). 
The deployment is composed of three diferent services:
| Docker alias          | Service Name       |
| ----------------------|--------------------|
| supervisorapi         | Supervisor API     |
| db                    | MySql Database     |
| queue                 | Redis              |
```yaml
version: '3'
services:
  supervisorapi:
    build:
      context: .
    image: mrbolas/supervisorapi:0.1.0
    ports:
      - "8080:8080" # http
    env_file:
      - .env-docker
    depends_on:
      - db
      - queue
    restart: on-failure

  db:
    image: mysql:5.7
    ports:
      - "3306:3306"
    environment:
      MYSQL_USER: user
      MYSQL_PASSWORD: password
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: sh_supervisor
    volumes:
      - db-data:/var/lib/mysql
    restart: unless-stopped

  queue:
    image: redis
    restart: always
    ports:
      - '6379:6379'

volumes:
 db-data:
```
### Run Docker Compose
```
docker-compose up
```
### Monitor Docker Compose Redis
To validate incoming notifications we can use:
```
docker exec -it supervisorapi_queue_1 redis-cli -h localhost subscribe notifications
```
This command subscribes to the topic notifications in redis-cli inside the docker container some-redis, and will show received messages received for that topic.
## K8s Environment

### K8s Requirements
1. Docker Destop running with K8s active
### K8s Environment variables
K8s deployment uses the same environment variables has [docker-compose](#docker-environment-variables)

### Setting Up K8s Environment
K8s files are generated by [kompose](https://kompose.io/) from the docker-compose file.
```
kompose --file docker-compose.yaml convert
```
These files are stored in the /k8s directory inside the project.
### Run K8s
To deploy the project in Kubernetes navigate to supervisorapi/k8s apply the K8s configuration files generated by Komposer. This will expose the supervisorAPI endpoints on port 30080
```
kubectl apply -f supervisorapi-pod.yaml,db-data-persistentvolumeclaim.yaml,db-deployment.yaml,db-service.yaml,queue-deployment.yaml,queue-service.yaml,supervisorapi-pod.yaml,supervisorapi-service.yaml
```
Additionally, it's needed to create a config map for the env-docker file wherre we store the environment variables.
```
kubectl create configmap env-docker --from-env-file=.env-docker
```

# Available Endpoints
(needs work)

# Testing and Coverage
(needs work)
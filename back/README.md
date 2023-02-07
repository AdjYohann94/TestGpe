## GPE Backend

This repository is the backend part of the GPE annual project.

It was a Golang web service that defined **REST** endpoints for the frontend application.

## Documentation

**Endpoints** : the documentation of all endpoints is available in [api/swagger.json](api%2Fswagger.json) file.

Once project is running : http://localhost:8081/

## Getting Started

### Prerequisites :

Taskfile : https://taskfile.dev

Copy and paste `env.example` into `.env`.

**Running the project with docker:**

- Docker and compose
- Go 1.19 (dev only)

**Running without docker:**

- Postgresql database
- Go 1.19 at least

### How to start (development)

**With docker :**

```shell
docker-compose up -d
```

On the first up, the database could not be ready at the time for migration job.
To rerun migrations use :

```shell
docker-compose restart migrate
```

This docker compose start the postgresql database, exposing openapi and runs migrations.
To start the go webserver application run :

```shell
go run cmd/app/main.go
```

or

```shell
task run
```

In this case, database must be accessible through localhost, default configuration is to expose at port `5432`.
Data is saved into a docker named volume.

Web service is available at :

### How to start (compiling)

Run with go Dockerfile :

```shell
docker-compose -f docker-compose-build.yml up -d
```

## Logging

Applications logs could be found in `log/debug.log` file. There is several log levels. The minimul log level could be
updated in `configs/main.yml` file.

## Configuration

Configuration is available in `.env` file for database connection credentials.
Used in docker-compose and loaded by godotenv.

There is some other configuration values in `main.yml` file for jwt lifetime.
These values are shipped into the containers on build.

## Commands

To see the logs :

```shell
task logs
```

To run the server :

```shell
task run
```

To build the binary :

```shell
task build-binary
```

To build golang docker image :

```shell
task build-image
```

Run tests :

```shell
task test
```

Run linters :

```shell
task lint
```

## Seeder :

View seed command line documentation :

```shell
task seed
```
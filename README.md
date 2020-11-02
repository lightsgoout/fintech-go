### API Documentation

See docs/api.md

### How to run the project

#### Serve API

`docker-compose up fintech`

API will be available at `http://localhost:8080/`

#### Test

`docker-compose up --force-recreate fintech_test`

Will run test suite against throwaway dockerized Postgres.

Current coverage is `84.5%`

#### Loadtest

`docker-compose up --force-recreate fintech_loadtest`

This will shoot various random requests to the API.

### Design overview

#### Project layout

```
payments/api - API for the service
payments/entity - business entities (Account, Payment)
payments/service - business logic interface
payments/service/persistent - business logic implementation based on Postgres
pkg/money - custom Money type (see rationale below)
pkg/postgres and pkg/testing - deal with postgres test isolation
```

#### Currency

Currency is stored as a Postgres ENUM as I felt like having a separate currency table would be an overkill 
for this task and storing currency as plain text seemed wrong.

Only `USD, EUR, RUB` currencies are supported.

#### Money

`pkg/money` contains custom `Numeric` type which is 
backed by popular `shopspring/decimal` lib for now, but we 
could easily switch to any other decimal lib.

In the database money stored as Postgres `numeric` data type.

#### Docker

Postgres image has a custom Dockerfile 
to work around MacOS security volume limitations (`init.sql is not shared from the host and is not known to Docker`) 
so you don't have to touch your docker configuration :)

### Topics out of scope of this task

Some things that make sense but were omitted for simplicity and to not bloat the project:

* Table partitioning. Payment and Account tables could easily be partitioned.
* Metrics/logging collection.
* Migrations. For now init.sql is copied into Postgres container.


 
# HR System

## Overview

HR System is a simple human resources management API Hub that provides functions such as User registration, login/logout, Employee information, clock in/out , Department management and processing leave requests.

## Features

- Employee
  - Register New User
  - Login/out User
  - Remove User
  - Update User's profiles
  - Reset User's password

- Department
  - CRUD Department

- Attendance
  - Submit and approve leave requests
  - Track leave history
  - Record employee chock-in/clock-out times

- Access Control
  - Role & Ability Model
  - User's Ability Authorization

## Technology Stack
- Backend:
  - [Go (Golang)](https://golang.org/): A fast, statically typed, compiled language
  - [Gin](https://github.com/gin-gonic/gin): A HTTP web framework
  - [GORM](https://gorm.io/): A ORM library for Golang
  - [fx](https://github.com/uber-go/fx): A dependency injection system for Go

- Database:
  - [MySQL](https://www.mysql.com/): Open-source relational database management system

- Caching:
  - [Redis](https://redis.io/): In-memory data structure store used as a database, cache, and message broker

- Authentication:
  - JWT (JSON Web Tokens) for secure authentication


## Project Structure

The project follows a structured layout for better organization. Below is an overview of the project structure:

```
/hr-system-go
  /app
    /plugins
    application.go
  /cmd
    /database
      main.go
    api.go
  /internal
    /user
  /database
  /utils
  go.mod
  go.sum
```

### app

Required tools for building apps, include mysql, http, redis. And use fx to Implement DI
- mysql: data store
- redis: cache store
- http: web framework
- env: load .env file
- logger: print logger logic

### cmd

Api Server & Database command(migration & seed) entrypoint

### internal

Implement application api
- user: Implement CRUD User API
- auth: Implement User's Role, Ability and Authorization logic
- attendance: Implement CRUD User's Leave and ClockIn/Out API
- department: Implement CRUD Department API
- session: Implement Register, login, logout User and resetPassword API

### database

Implement mysql database migration & seed logic with GORM

## Getting Started

1. Clone repo

```
git clone https://github.com/kken339039/hr-system-go
cd hr-system-go
```

2. Build and run the Docker container:

When running for the first time, you need to wait for Docker Compose to prepare the database and Redis before the API server can successfully connect.

This make command include build & run in binary

```
make start
```

## Database

Database Migration & Seed files and implementation logic for operating DB, use the following command:

- Initialize Database
```
make db-init
```

- Create new migration file
```
make db-migration-create ${fileName}
```

- Run migrations those has not yet been run
```
make db-migration-run
```

- Rollback last migration
```
make db-seed-create
```

- Create New migration file
```
make db-seed-create ${fileName}
```

- Run all seed files
```
make db-seed-run-all
```

- Run specific seed files
```
make db-seed-run ${fileName}
```

## Local Development
### Development Tool

- [golangci-lint](https://github.com/golangci/golangci-lint): A runner for Go linters
  - Used for code formatting and maintaining code quality

- [Air](https://github.com/cosmtrek/air): Live reload for Go apps
  - Enables hot reloading during development

### Setup the HTTP server in local

1. install Go 1.22.5
2. copy `.env.example` as `.env` and update the values(ex: database host -> 127.0.0.1)
3. Run `go mod download` to download the go modules
4. go install github.com/cosmtrek/air@latest
5. Run `air` to start the HTTP server
6. Use postman or to call api and start development


## Format & Lint Code

To format & lint code follow golangci.yml, then you can fix it if any warning or error. use the following command:

```
make format
make lint
```

## Unit Tests

To run unit tests, use the following command:

```
make test
```

## TBD
### CORS
- Using CORS which gin support

### Mailer
- Implement Mailer for send reset password token or some reports

### API Doc
- Generate Doc by Swagger

### Role & Ability API
- CRUD New Role & Ability API

### Monitoring and Logging:
- Implement monitor and logging solutions to track system error and diagnose issues.

### CI/CD Pipeline:
- Set up a continuous integration/continuous deployment (CI/CD) pipeline for automated testing and deployment.
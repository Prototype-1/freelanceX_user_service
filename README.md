# User Service — FreelanceX

## Overview
User service focus more on authentication/authorization using JWT for users (admin, freelancers and clients). We use Postgres as the database. This service beholds several real business oriented constraints.

## Tech Stack
- Go (Golang)
- JWT
- gRPC
- PostgreSQL
- Protocol Buffers
- Redis (optional: for session or token caching)

## Folder Structure
.
├── config/ # Configuration and environment loading
├── handler/ # gRPC handler implementations
├── model/ # Data models for DB
├── proto/ # Protobuf service definitions
├── repository/ # DB interactions via GORM
├── service/ # Business logic layer
├── main.go # App entrypoint
└── go.mod


## Setup

### 1. Clone & Navigate
```bash
git clone https://github.com/Prototype-1/freelancex_user_service.git
cd freelancex_project/user_service
```

### Install Dependencies

go mod tidy

## .env

PORT=50051
DB_URL=postgres://username:password@localhost:5432/userdb
JWT_SECRET=your_jwt_secret
REDIS_ADDR=localhost:6379

### Run Migrations

go run scripts/migrate.go

##  Start the Service

go run main.go

## Proto Definitions

Located in proto/auth/, proto/profile/, proto/portfolio/, and proto/review/.

To regenerate:

protoc --go_out=. --go-grpc_out=. proto/auth/user.proto
protoc --go_out=. --go-grpc_out=. proto/auth/review.proto
protoc --go_out=. --go-grpc_out=. proto/auth/profile.proto
protoc --go_out=. --go-grpc_out=. proto/auth/portfolio.proto

### Ports

    gRPC: :50051

#### Notes

    This service auto-creates empty profiles for freelancers on OAuth.

    Only one admin is supported system-wide.
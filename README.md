# finance-dashboard-api

## Overview

A backend system for managing financial records, users, and role-based access control, built with a focus on clean architecture, data integrity, and scalable API design.

This project demonstrates backend development skills including API design, business logic implementation, access control, and data modeling.

## Tech Stack

- Language: Go
- Router: Chi
- Database: PostgreSQL
- ORM/Query Tool: SQLC
- Migrations: Goose
- Authentication: JWT

## Features

### User & Role Management

- Create and manage users
- Assign roles (Viewer, Analyst, Admin)
- Enable/disable users
- Role-based access control enforcement

### Financial Records Management

- Create, read, update, delete records
- Fields include:
- Amount
- Type (income/expense)
- Category
- Date
- Notes
- Filtering support (type, category, date)

### Dashboard & Analytics APIs

- Total income
- Total expenses
- Net balance
- Category-wise breakdown
- Monthly trends
- Recent activity

### Access Control

- Viewer → Read-only access
- Analyst → Read + insights
- Admin → Full control

### Validation & Error Handling

- Input validation
- Structured error responses
- Proper HTTP status codes
- Safe handling of invalid operations

### Data Persistence

- PostgreSQL database with structured schema
- Efficient queries using SQLC

## Design Decisions

- Used SQLC instead of ORM for better type safety and performance
- Implemented role based access control using middleware for clean separation
- Chose PostgreSQL for relational consistency and aggregation queries
- Used soft deletes to preserve historical data
- Structured handlers to keep logic simple and readable for this scope

## Project Structure

```
.
├── main.go                 # Application entry point
├── go.mod / go.sum        # Dependencies

├── app/                   # App setup & routing
│   ├── app.go             # App initialization
│   └── routes.go          # Route definitions

├── handler/               # HTTP handlers (controllers)
│   ├── handler.go
│   ├── handler_auth.go
│   ├── handler_user.go
│   ├── handler_record.go
│   ├── handler_dashboard.go
│   └── middleware.go      # Auth & role-based middleware

├── internal/
│   └── database/          # SQLC generated DB layer
│       ├── db.go
│       ├── models.go
│       ├── users.sql.go
│       ├── records.sql.go
│       ├── dashboard.sql.go
│       └── sessions.sql.go

├── models/                # Request/response models (DTOs)
│   ├── user.go
│   ├── record.go
│   └── dashboard.go

├── sql/                   # SQL source files
│   ├── schema/            # Database schema & migrations
│   │   ├── 001_users.sql
│   │   ├── 002_records.sql
│   │   ├── 003_sessions.sql
│   │   ├── 004_created_by_records.sql
│   │   └── 005_indexes.sql
│   └── queries/           # SQL queries for SQLC
│       ├── users.sql
│       ├── records.sql
│       ├── dashboard.sql
│       └── sessions.sql

├── token/                 # JWT handling & roles
│   ├── jwt_maker.go
│   ├── claims.go
│   └── roles.go

├── utils/                 # Utility functions
│   ├── helper.go
│   ├── json.go
│   └── password.go

├── validators/            # Input validation logic
│   └── validate_record.go

├── seeds/                 # Seed data for testing
│   ├── records.json
│   └── seed.js

├── api-docs/                   # API collections (Postman)
├── sqlc.yaml              # SQLC configuration
└── README.md
```

## API Overview

### Base Routes

| Method | Endpoint | Description         | Access |
| ------ | -------- | ------------------- | ------ |
| GET    | `/`      | Health check        | Public |
| GET    | `/err`   | Test error response | Public |

---

## Auth & User Routes (`/users`)

### Public

| Method | Endpoint              | Description        |
| ------ | --------------------- | ------------------ |
| POST   | `/users/login`        | Login user         |
| POST   | `/users/tokens/renew` | Renew access token |

### Protected

| Method | Endpoint               | Description       | Role          |
| ------ | ---------------------- | ----------------- | ------------- |
| GET    | `/users/logout`        | Logout user       | Authenticated |
| POST   | `/users/tokens/revoke` | Revoke session    | Authenticated |
| GET    | `/users/deleted`       | Get deleted users | Admin         |
| POST   | `/users/`              | Create user       | Admin         |
| GET    | `/users/`              | List users        | Admin         |
| DELETE | `/users/{id}`          | Soft delete user  | Admin         |
| DELETE | `/users/{id}/h`        | Hard delete user  | Admin         |

---

## Record Routes (`/records`)

| Method | Endpoint           | Description         | Role    |
| ------ | ------------------ | ------------------- | ------- |
| GET    | `/records/{id}`    | Get record by ID    | Analyst |
| GET    | `/records/`        | List records        | Analyst |
| POST   | `/records/`        | Create record       | Admin   |
| PATCH  | `/records/{id}`    | Update record       | Admin   |
| DELETE | `/records/{id}`    | Soft delete record  | Admin   |
| GET    | `/records/deleted` | Get deleted records | Admin   |
| DELETE | `/records/{id}/h`  | Hard delete record  | Admin   |

---

## Dashboard Routes (`/dashboard`)

| Method | Endpoint                | Description                | Role   |
| ------ | ----------------------- | -------------------------- | ------ |
| GET    | `/dashboard/summary`    | Get income/expense summary | Viewer |
| GET    | `/dashboard/categories` | Category-wise analysis     | Viewer |
| GET    | `/dashboard/trends`     | Monthly trends             | Viewer |
| GET    | `/dashboard/recent`     | Recent transactions        | Viewer |

---

## Access Control Summary

- **Viewer** → Can access dashboard endpoints
- **Analyst** → Can view records + dashboard
- **Admin** → Full access (users + records management)

---

## Prerequisites

Make sure you have the following installed before running the project

### Golang

- Install from: <https://go.dev/doc/install>
- Verify Installation:

```bash
go version
```

### SQLC

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

- Verify Installation

```bash
sqlc version
```

### Goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

- Verify Installation

```bash
goose -version
```

## Setup Instructions

### Clone the repo

```bash
git clone https://github.com/aditya-sutar-45/finance-dashboard-api.git
cd finance-dashboard-api
```

### Setup environment variables

```env
PORT=:3000
DB_URL=<YOUR_DB_URL>?sslmode=disable
SECRET_KEY=32948903840981234
```

### Run Database Migrations

```bash
cd sql/schema
goose <YOUR_DB_URL> up
cd ../../
```

### Generate SQLC code

```bash
sqlc generate
```

### Build and Run the server

```bash
go mod tidy
go build -o api main.go && ./api
```

## Optional Enhancements Implemented

- Pagination for records
- Filtering & query params
- Soft delete support
- API testing via Postman

## Assumptions

- Users must be authenticated for most endpoints
- Admin has full control over system
- Records are scoped per user (unless admin)
- Dates are handled in UTC

## Future Improvements

- Unit & integration tests
- Rate limiting
- Caching for dashboard endpoints
- WebSocket-based real-time updates
- Frontend dashboard integration
- More robust Analysis layer

## Author

Aditya Girish Sutar

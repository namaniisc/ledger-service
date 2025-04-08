# Ledger Service

A simple RESTful ledger service in Go, using Gin and PostgreSQL.  
Supports creating customers, credit/debit transactions (with idempotency), and retrieving balances & transaction history.

---

## Table of Contents

- [Features](#features)  
- [Prerequisites](#prerequisites)  
- [Getting Started](#getting-started)  
  - [1. Clone Repository](#1-clone-repository)  
  - [2. Configure Environment](#2-configure-environment)  
  - [3. Install Dependencies](#3-install-dependencies)  
  - [4. Run Database Migrations](#4-run-database-migrations)  
  - [5. Start the Service](#5-start-the-service)  
- [API Endpoints](#api-endpoints)  
- [Development](#development)  
- [Project Structure](#project-structure)  

---

## Features

- **Customer Management**  
  - Create a new customer with initial balance  
  - Retrieve current balance  

- **Transactions**  
  - Credit or debit a customer account  
  - Idempotency via optional `client_transaction_id`  
  - Retrieve transaction history  

- **Database Migrations**  
  - SQL files under `migrations/` for schema setup  

---

## Prerequisites

- Go 1.18+  
- PostgreSQL 12+  
- `git`  

---

## Getting Started

### 1. Clone Repository

```bash
git clone https://github.com/namaniisc/ledger-service.git
cd ledger-service
```

### 2. Configure Environment

Copy the example and fill in your values:

```bash
cp .env.example .env
```

Example `.env` file:

```
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ledgerdb
```

> **Security:** `.env` is listed in `.gitignore`, so your secrets stay local.

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run Database Migrations

```bash
psql -U $DB_USER -h $DB_HOST -p $DB_PORT -d $DB_NAME -f migrations/001_init.up.sql
```

### 5. Start the Service

```bash
go run main.go
```

---

## API Endpoints

| Method | Path                                   | Description                             |
|--------|----------------------------------------|-----------------------------------------|
| POST   | `/customers`                           | Create a new customer                   |
| GET    | `/customers/:customer_id/balance`      | Get current balance for a customer      |
| POST   | `/transactions`                        | Create a credit/debit transaction       |
| GET    | `/customers/:customer_id/transactions` | List all transactions for a customer    |

---

### Create Customer

**POST**  
`http://localhost:8080/customers`

**Request**
```json
{
  "name": "Naman",
  "initial_balance": 1000
}
```

**Response**
```json
{
  "customer_id": "ab274ae1-10ce-4478-9c67-2c0a8ebfdc26",
  "name": "Naman",
  "balance": 1000
}
```

---

### Credit Transaction

**POST**  
`http://localhost:8080/transactions`

**Request**
```json
{
  "customer_id": "ab274ae1-10ce-4478-9c67-2c0a8ebfdc26",
  "type": "credit",
  "amount": 500
}
```

**Response**
```json
{
  "balance": 1500,
  "status": "success",
  "transaction_id": "1e621224-9129-462b-a068-ad82ff85e06d"
}
```

---

### Debit Transaction

**POST**  
`http://localhost:8080/transactions`

**Request**
```json
{
  "customer_id": "ab274ae1-10ce-4478-9c67-2c0a8ebfdc26",
  "type": "debit",
  "amount": 100
}
```

**Response**
```json
{
  "balance": 1400,
  "status": "success",
  "transaction_id": "820603ca-779d-4fa1-b9c3-7ee067004309"
}
```

---

### Get Balance

**GET**  
`http://localhost:8080/customers/ab274ae1-10ce-4478-9c67-2c0a8ebfdc26/balance`

**Response**
```json
{
  "balance": 1400,
  "customer_id": "ab274ae1-10ce-4478-9c67-2c0a8ebfdc26"
}
```

---

### Get All Transactions

**GET**  
`http://localhost:8080/customers/ab274ae1-10ce-4478-9c67-2c0a8ebfdc26/transactions`

**Response**
```json
[
  {
    "transaction_id": "820603ca-779d-4fa1-b9c3-7ee067004309",
    "customer_id": "00000000-0000-0000-0000-000000000000",
    "type": "debit",
    "amount": 100,
    "timestamp": "2025-04-08T20:31:47.889028Z"
  },
  {
    "transaction_id": "1e621224-9129-462b-a068-ad82ff85e06d",
    "customer_id": "00000000-0000-0000-0000-000000000000",
    "type": "credit",
    "amount": 500,
    "timestamp": "2025-04-08T20:31:00.556717Z"
  }
]
```

---

## Development

- **Environment Variables:** Loaded via [`joho/godotenv`](https://github.com/joho/godotenv)  
- **Routing & HTTP:** [`gin-gonic/gin`](https://github.com/gin-gonic/gin)  
- **Postgres Driver:** [`lib/pq`](https://github.com/lib/pq)  

```bash
go test ./...
```

---

## Project Structure

```
.
├── .env.example        # Example env file (no secrets)
├── .gitignore
├── go.mod
├── go.sum
├── main.go             # Application entrypoint
├── database/
│   └── database.go     # DB initialization
├── handlers/
│   ├── customer.go     # Customer endpoints
│   └── transaction.go  # Transaction endpoints
├── models/
│   └── model.go        # Data models
└── migrations/
    └── 001_init.up.sql # Initial schema
```

---


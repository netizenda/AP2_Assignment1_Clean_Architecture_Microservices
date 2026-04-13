# AP2 Assignment 2: Clean Architecture Microservices with gRPC

## Project Description

This project demonstrates **Clean Architecture** in a microservices environment using **Go**, **Gin** (REST) and **gRPC** (internal communication). It follows the **Contract-First** approach with separate repositories for `.proto` definitions and generated code.

- **Order Service** — REST API + gRPC streaming server
- **Payment Service** — gRPC server
- Separate PostgreSQL databases
- Full business logic in Use Cases layer

## Architecture

- **Domain** — business entities
- **Usecase** — business rules and orchestration
- **Repository** — data access (PostgreSQL)
- **Transport** — HTTP (Gin) + gRPC
- **Client** — gRPC client for inter-service communication

## Services

| Service         | Transport          | Port   | Database     |
|-----------------|--------------------|--------|--------------|
| Order Service   | REST + gRPC Stream | 8080 / 50052 | `orderdb`    |
| Payment Service | gRPC               | 50051  | `paymentdb`  |

## Technologies

- Go 1.25
- Gin (REST)
- gRPC + Protocol Buffers
- PostgreSQL
- Clean Architecture
- GitHub Actions (Contract-First code generation)

## Project Structure

AP2_Assignment1_Clean_Architecture_Microservices/
├── order-service/
│   ├── cmd/order-service/main.go
│   ├── internal/
│   │   ├── app/
│   │   ├── client/
│   │   ├── domain/
│   │   ├── repository/
│   │   ├── transport/http/
│   │   ├── transport/grpc/
│   │   └── usecase/
│   └── proto/v1/                  ← local copy of generated proto
├── payment-service/
│   ├── cmd/payment-service/main.go
│   └── ... (similar structure)
├── AP2_Protos/                     ← source .proto (separate repo)
└── AP2_Generated/                  ← generated code (separate repo)


## How to Run

1. Start PostgreSQL and create two databases:
    - `orderdb`
    - `paymentdb`

2. Run migrations (files in `migrations/` folder).

3. Start **Payment Service** first:
   ```powershell
   cd payment-service
   go run cmd/payment-service/main.go
4. Start **Order Service**
   ```powershell
    cd order-service
    go run cmd/order-service/main.go

## Test Requests
    ```powershell
    $body = @{ customer_id = "cust-001"; item_name = "iPhone 16"; amount = 15000 } | ConvertTo-Json

    Invoke-WebRequest -Method POST `
    -Uri "http://localhost:8080/orders" `
    -Headers @{ "Content-Type" = "application/json" } `
    -Body $body

## Contract-First gRPC

.proto files are stored in separate repo AP2_Protos

Code is auto-generated via GitHub Actions into AP2_Generated

Payment uses unary RPC

Order uses server-side streaming
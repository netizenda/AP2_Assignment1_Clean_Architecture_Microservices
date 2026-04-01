# AP2 Assignment 1 – Clean Architecture based Microservices (Order & Payment)

**Student:** Taubakabyl Nurlybek  
**Group:** [Write your group here]  
**Date:** April 01, 2026

## Project Overview

This project implements a two-service platform (Order Service + Payment Service) following **Clean Architecture** principles and microservices best practices. Communication is done strictly via REST using Gin framework.

## Architecture

Each service is built according to **Clean Architecture**:
- `domain/` – Pure business entities (independent of frameworks and databases)
- `usecase/` – All business logic and invariants
- `repository/` – Data access layer with PostgreSQL (interface + implementation)
- `transport/http/` – Thin delivery layer (only request/response handling)
- `client/` (Order Service only) – Outbound HTTP client with 2-second timeout
- `app/` – Composition Root with manual dependency injection

**Strict Microservices Rules Applied:**
- Separate databases (`orderdb` and `paymentdb`)
- No shared code or models between services
- Synchronous REST communication only

## Bounded Contexts

- **Order Service** – Responsible for order lifecycle and state management
- **Payment Service** – Responsible for payment processing and transaction limits

## How to Run

```powershell
# Terminal 1 - Payment Service
cd payment-service
$env:DATABASE_URL = "postgres://postgres:12345@localhost:5432/paymentdb?sslmode=disable"
$env:PORT = "8081"
go run cmd/payment-service/main.go

# Terminal 2 - Order Service
cd order-service
$env:DATABASE_URL = "postgres://postgres:12345@localhost:5432/orderdb?sslmode=disable"
$env:PAYMENT_SERVICE_URL = "http://localhost:8081"
$env:PORT = "8080"
go run cmd/order-service/main.go

API Testing Examples (PowerShell)
# 1. Create Order (with Idempotency)
$body = @{customer_id="cust-001"; item_name="iPhone 16"; amount=15000} | ConvertTo-Json
Invoke-WebRequest -Method POST -Uri "http://localhost:8080/orders" `
    -Headers @{"Content-Type"="application/json"; "X-Idempotency-Key"="key-123"} -Body $body

# 2. Get Order
Invoke-WebRequest -Uri "http://localhost:8080/orders/{order_id}"

# 3. Cancel Order (only works for Pending status)
Invoke-WebRequest -Method PATCH -Uri "http://localhost:8080/orders/{order_id}/cancel"

Bonus Feature Implemented
Idempotency – using X-Idempotency-Key header to prevent duplicate orders/payments.

Failure Scenario (Required by Assignment)
When Payment Service is unavailable:

Order Service uses custom http.Client with 2-second timeout
Returns 503 Service Unavailable
Order is marked as "Failed"
```
## Architecture Diagram
![img.png](img.png)

## Key Design Decisions
Full separation of concerns and dependency inversion

Business rules are located only in the use case layer

Domain models do not depend on HTTP, JSON, or database structures

Each service owns its own data completely

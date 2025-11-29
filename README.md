
# E-Commerce API

A simple e-commerce backend API built with Go. Supports products management and order creation with transactional handling.

## Features

* List, create, update, and delete products
* Place orders with multiple items
* Update product stock automatically on order creation
* Transactional order creation to ensure data consistency
* Healthcheck endpoint

## Setup

1. Clone the repo:

```bash
git clone git@github.com:Stephen-Kamau/go-ecommerce-api.git
cd go-ecommerce-api
```

2. Configure the database in `config`:

```env
DATABASE_URL=postgres://user:password@localhost:5432/ecomdb?sslmode=disable
APP_ADDRESS=:8080
```

3. Run migrations with Goose:

```bash
goose -dir migrations postgres "$DATABASE_URL" up
```

4. Generate SQL query code with sqlc:

```bash
sqlc generate
```

5. Run the server:

```bash
go run main.go
```

Server starts at `http://localhost:8080`.

## API Endpoints

### Products

| Method | Path           | Description          |
| ------ | -------------- | -------------------- |
| GET    | /products      | List all products    |
| POST   | /products      | Create a new product |
| GET    | /products/{id} | Get product by ID    |
| DELETE | /products/{id} | Delete product       |

### Orders

| Method | Path                   | Description                |
| ------ | ---------------------- | -------------------------- |
| POST   | /orders                | Place a new order          |
| GET    | /orders                | Get all orders             |
| GET    | /orders/{id}           | Get order by ID            |
| GET    | /orders/customer/{ref} | Get orders by customer ref |

### Healthcheck

| Method | Path    | Description      |
| ------ | ------- | ---------------- |
| GET    | /health | Check API status |


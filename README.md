# InvoiceGen

A simple invoice generator built with Go, MySQL, and plain HTML/CSS/JS. Create invoices, manage customers, and download PDF copies.

## Features

- Create invoices with multiple line items
- Auto-fill customer details by phone number
- Save and manage customers
- Generate and download invoices as PDF
- View all generated invoices

## Tech Stack

- **Backend** — Go (standard library `net/http`)
- **Database** — MySQL 8.0 via `sqlx`
- **PDF** — `github.com/go-pdf/fpdf`
- **Frontend** — HTML, CSS, Vanilla JS
- **Infrastructure** — Docker Compose

## Project Structure

```
Invoice Generator/
├── cmd/
│   └── main.go                 # Entry point, HTTP server, routes
├── internal/
│   ├── db/
│   │   └── db.go               # MySQL connection
│   ├── handlers/
│   │   └── web.go              # HTTP handlers
│   ├── models/                 # Data models
│   ├── pdf/
│   │   └── generator.go        # PDF generation
│   └── repository/
│       ├── customer.go         # Customer DB queries
│       └── invoice.go          # Invoice DB queries
├── migrations/
│   └── 001_init.sql            # Database schema
├── web/
│   ├── static/
│   │   ├── css/style.css
│   │   ├── js/app.js
│   │   └── favicon.svg
│   └── templates/              # Go HTML templates
├── Dockerfile
└── docker-compose.yml
```

## Getting Started

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Docker](https://www.docker.com/)

### 1. Start MySQL

```bash
docker compose up db -d
```

This starts MySQL 8.0 on port `3306` and auto-runs the migration.

### 2. Run the server

```bash
go run ./cmd/main.go
```

Open [http://localhost:8080](http://localhost:8080)

### Run with Docker (app + db)

```bash
docker compose up --build
```

## Routes

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/` | Create invoice form |
| `POST` | `/invoice/generate` | Save invoice and redirect to view |
| `GET` | `/invoice/{id}` | View invoice |
| `GET` | `/invoice/{id}/pdf` | Download invoice as PDF |
| `GET` | `/invoices` | List all invoices |
| `GET` | `/customers` | Customer form + saved customers |
| `POST` | `/customers/save` | Save a customer |
| `GET` | `/api/customer?phone=` | Lookup customer by phone (JSON) |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MYSQL_DSN` | `root:root@tcp(127.0.0.1:3306)/invoice?parseTime=true` | MySQL connection string |
| `PORT` | `8080` | HTTP server port |

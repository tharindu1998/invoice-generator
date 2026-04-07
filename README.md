# InvoiceGen

A web-based invoice generator built with Go, MySQL, and plain HTML/CSS/JS. Create professional invoices, manage customers, and download PDF copies.

## Features

- Create invoices with seller & customer details (name, address, phone)
- Auto-fill customer details by phone number lookup
- Save and manage customers
- Generate and download invoices as PDF
- View all generated invoices with delete support

## Tech Stack

- **Backend** — Go (standard library `net/http`)
- **Database** — MySQL 8.0 via `sqlx`
- **PDF** — `github.com/go-pdf/fpdf`
- **Frontend** — HTML, CSS, Vanilla JS
- **Infrastructure** — Docker Compose
- **CI/CD** — GitHub Actions → Render deploy hook

## Project Structure

```
Invoice Generator/
├── cmd/
│   └── main.go                 # Entry point, HTTP server, routes
├── internal/
│   ├── db/
│   │   └── db.go               # MySQL connection (reads MYSQL_DSN env)
│   ├── handlers/
│   │   └── web.go              # HTTP handlers
│   ├── models/                 # Data models (Invoice, Customer, PaymentInfo, etc.)
│   ├── pdf/
│   │   └── generator.go        # PDF generation
│   └── repository/
│       ├── customer.go         # Customer DB queries
│       └── invoice.go          # Invoice DB queries
├── migrations/
│   └── 001_init.sql            # Full database schema (run once)
├── web/
│   ├── static/
│   │   ├── css/style.css
│   │   ├── js/app.js
│   │   ├── logo.svg
│   │   └── favicon.svg
│   └── templates/              # Go HTML templates
│       ├── base.html
│       ├── create-invoice.html
│       ├── customer-form.html
│       ├── invoice-preview.html
│       └── invoices.html
├── .github/
│   └── workflows/
│       └── render-deploy.yml   # Auto-deploy to Render on push to main
├── Dockerfile
└── docker-compose.yml
```

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/)

### Run with Docker

```bash
docker compose up --build
```

Open [http://localhost:8080](http://localhost:8080)

This starts MySQL 8.0 and the app together. The migration in `migrations/001_init.sql` runs automatically on first start.

## Routes

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/` | Create invoice form |
| `POST` | `/invoice/generate` | Save invoice and redirect to view |
| `GET` | `/invoice/{id}` | View invoice |
| `GET` | `/invoice/{id}/pdf` | Download invoice as PDF |
| `POST` | `/invoice/{id}/delete` | Delete invoice |
| `GET` | `/invoices` | List all invoices |
| `GET` | `/customers` | Customer form + saved customers list |
| `POST` | `/customers/save` | Save or update a customer |
| `POST` | `/customers/{id}/delete` | Delete a customer |
| `GET` | `/api/customer?phone=` | Lookup customer by phone (JSON) |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MYSQL_DSN` | `root:root@tcp(127.0.0.1:3306)/invoice?parseTime=true` | MySQL connection string |
| `PORT` | `8080` | HTTP server port |

## DEMO

Demo Link : https://invoice-generator-2rzm.onrender.com/

<p align="center">
  <h3>Create Invoice View</h3>
  <img width="600" alt="Create Invoice View" src="https://github.com/user-attachments/assets/9279d17c-faf6-4efb-8f0c-a5fe50dabb71" />
</p>

<p align="center">
  <h3>Create Customer View</h3>
  <img width="600" height="880" alt="Create Customer View" src="https://github.com/user-attachments/assets/8a556506-0b64-42d4-892e-eb76a006cbfe" />
</p>

<p align="center">
  <h3>Preview Created Invoice</h3>
  <img width="600" height="958" alt="Preview Created Invoice" src="https://github.com/user-attachments/assets/c4333fef-b4b1-4249-89d1-3afdcf64bced" />
</p>

<p align="center">
  <h3>Downloaded Invoice PDF</h3>
  <img width="600" height="899" alt="Invoice" src="https://github.com/user-attachments/assets/c9f19b5d-e479-4d47-9773-1aab135545d7" />
</p>






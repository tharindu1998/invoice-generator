package main

import (
	"log"
	"net/http"
	"os"

	"invoice-generator/internal/db"
	"invoice-generator/internal/handlers"
)

func main() {
	if _, err := db.InitDB(); err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	log.Println("database connected")

	// Static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Pages
	http.HandleFunc("GET /{$}", handlers.CreateInvoice)
	http.HandleFunc("GET /customers", handlers.Customers)
	http.HandleFunc("GET /invoices", handlers.Invoices)
	http.HandleFunc("GET /invoice/{id}", handlers.ViewInvoice)
	http.HandleFunc("GET /invoice/{id}/pdf", handlers.DownloadPDF)

	// Actions
	http.HandleFunc("POST /customers/save", handlers.SaveCustomer)
	http.HandleFunc("POST /invoice/generate", handlers.GenerateInvoice)

	// API
	http.HandleFunc("GET /api/customer", handlers.LookupCustomer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server running → http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

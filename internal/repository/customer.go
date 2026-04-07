package repository

import (
	"database/sql"
	"errors"

	"invoice-generator/internal/db"
	"invoice-generator/internal/models"
)

var ErrNotFound = errors.New("not found")

func SaveCustomer(c models.Customer) error {
	_, err := db.DB.NamedExec(`
		INSERT INTO customers (name, email, address_line1, address_line2, phone)
		VALUES (:name, :email, :address_line1, :address_line2, :phone)
		ON DUPLICATE KEY UPDATE
			name          = VALUES(name),
			email         = VALUES(email),
			address_line1 = IF(VALUES(address_line1) != '', VALUES(address_line1), address_line1),
			address_line2 = IF(VALUES(address_line2) != '', VALUES(address_line2), address_line2)
	`, c)
	return err
}

func GetAllCustomers() ([]models.Customer, error) {
	var customers []models.Customer
	err := db.DB.Select(&customers, "SELECT * FROM customers ORDER BY created_at DESC")
	return customers, err
}

func DeleteCustomer(id int64) error {
	_, err := db.DB.Exec("DELETE FROM customers WHERE id = ?", id)
	return err
}

func GetCustomerByPhone(phone string) (models.Customer, error) {
	var c models.Customer
	err := db.DB.Get(&c, "SELECT * FROM customers WHERE phone = ? LIMIT 1", phone)
	if errors.Is(err, sql.ErrNoRows) {
		return c, ErrNotFound
	}
	return c, err
}

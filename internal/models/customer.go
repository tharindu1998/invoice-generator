package models

type Customer struct {
	Mobile       string `db:"mobile"`
	Name         string `db:"name"`
	Email        string `db:"email"`
	AddressLine1 string `db:"address_line1"`
	AddressLine2 string `db:"address_line2"`
}

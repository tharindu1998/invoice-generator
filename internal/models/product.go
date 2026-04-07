package models

type Product struct {
	ID          int64   `db:"id"          json:"id"`
	Name        string  `db:"name"        json:"name"`
	Description string  `db:"description" json:"description"`
	Price       float64 `db:"price"       json:"price"`
}

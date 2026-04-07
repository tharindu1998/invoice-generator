package models

import "time"

type Customer struct {
	ID           int64     `db:"id"           json:"id"`
	Phone        string    `db:"phone"        json:"phone"`
	Name         string    `db:"name"         json:"name"`
	Email        string    `db:"email"        json:"email"`
	AddressLine1 string    `db:"address_line1" json:"address_line1"`
	AddressLine2 string    `db:"address_line2" json:"address_line2"`
	CreatedAt    time.Time `db:"created_at"   json:"created_at"`
}

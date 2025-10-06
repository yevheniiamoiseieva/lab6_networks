package models

import "time"

type Employee struct {
	ID         int       `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Position   string    `db:"position" json:"position"`
	Department string    `db:"department" json:"department"`
	Salary     float64   `db:"salary" json:"salary"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

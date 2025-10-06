package repositories

import "github.com/jmoiron/sqlx"

type Repositories struct {
	EmployeeRepository *EmployeeRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		EmployeeRepository: NewEmployeeRepository(db),
	}
}

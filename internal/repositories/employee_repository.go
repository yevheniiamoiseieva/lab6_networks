package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"laba6/internal/models"
)

type EmployeeRepository struct {
	db *sqlx.DB
}

func NewEmployeeRepository(db *sqlx.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) GetAll() ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.Select(&employees, "SELECT * FROM employees ORDER BY id")
	return employees, err
}

func (r *EmployeeRepository) Create(name, position, department string, salary float64) error {
	// Validate salary
	if salary < 0 {
		return fmt.Errorf("salary cannot be negative")
	}

	// Validate required fields
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if position == "" {
		return fmt.Errorf("position is required")
	}
	if department == "" {
		return fmt.Errorf("department is required")
	}

	// Check for duplicate name
	var count int
	checkQuery := `SELECT COUNT(*) FROM employees WHERE name = $1`
	err := r.db.Get(&count, checkQuery, name)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("employee with name '%s' already exists", name)
	}

	// Create employee
	query := `
       INSERT INTO employees (name, position, department, salary)
       VALUES ($1, $2, $3, $4)
    `
	_, err = r.db.Exec(query, name, position, department, salary)
	return err
}
func (r *EmployeeRepository) GetByID(id int) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.Get(&employee, "SELECT * FROM employees WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *EmployeeRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM employees WHERE id=$1", id)
	return err
}

func (r *EmployeeRepository) Update(id int, name, position, department string, salary float64) error {
	query := `
        UPDATE employees 
        SET name=$1, position=$2, department=$3, salary=$4 
        WHERE id=$5
    `
	result, err := r.db.Exec(query, name, position, department, salary, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("employee not found")
	}

	return nil
}

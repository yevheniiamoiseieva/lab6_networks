package processors

import (
	"laba6/internal/models"
	"laba6/internal/repositories"
)

type EmployeeProcessor struct {
	repo *repositories.EmployeeRepository
}

func NewEmployeeProcessor(repo *repositories.EmployeeRepository) *EmployeeProcessor {
	return &EmployeeProcessor{repo: repo}
}

func (p *EmployeeProcessor) GetAllEmployees() ([]models.Employee, error) {
	return p.repo.GetAll()
}

func (p *EmployeeProcessor) GetEmployeeByID(id int) (*models.Employee, error) {
	return p.repo.GetByID(id)
}

func (p *EmployeeProcessor) CreateEmployee(name, position, department string, salary float64) error {
	return p.repo.Create(name, position, department, salary)
}

func (p *EmployeeProcessor) DeleteEmployee(id int) error {
	return p.repo.Delete(id)
}

func (p *EmployeeProcessor) UpdateEmployee(id int, name, position, department string, salary float64) error {
	return p.repo.Update(id, name, position, department, salary)
}

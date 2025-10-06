package processors

import "laba6/internal/repositories"

type Processors struct {
	EmployeeProcessor *EmployeeProcessor
}

func NewProcessors(repos *repositories.Repositories) *Processors {
	return &Processors{
		EmployeeProcessor: NewEmployeeProcessor(repos.EmployeeRepository),
	}
}

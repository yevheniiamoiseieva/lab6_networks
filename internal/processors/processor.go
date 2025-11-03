package processors

import (
	"laba6/internal/repositories"
)

type Processors struct {
	EmployeeProcessor *EmployeeProcessor
	Rsa               IRsaService
	Aes               IAesService
}

func NewProcessors(repos *repositories.Repositories, rsaBits int, aesKeySize int) *Processors {
	return &Processors{
		EmployeeProcessor: NewEmployeeProcessor(repos.EmployeeRepository),
		Rsa:               NewRsaService(rsaBits),
		Aes:               NewAesService(aesKeySize),
	}
}

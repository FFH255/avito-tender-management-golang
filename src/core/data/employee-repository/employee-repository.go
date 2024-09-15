package employee_repository

import (
	"tms/src/core/services/repositories"
	"tms/src/pkg/pg"
)

type EmployeeRepository struct {
	client pg.Client
}

func New(client pg.Client) repositories.EmployeeRepository {
	return EmployeeRepository{client}
}

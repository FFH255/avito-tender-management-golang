package repositories

import (
	"context"
	"tms/src/core/domain"
)

type GetEmployeeDTO struct {
	ID       *domain.ID
	Username *string
}

type EmployeeRepository interface {
	Get(ctx context.Context, dto GetEmployeeDTO) (*domain.Employee, error)
}

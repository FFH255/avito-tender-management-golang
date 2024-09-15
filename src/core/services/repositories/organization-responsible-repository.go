package repositories

import (
	"context"
	"tms/src/core/domain"
)

type GetOrganizationResponsiblesListDTO struct {
	EmployeeID     *domain.ID
	OrganizationID *domain.ID
}

type GetOrganizationResponsibleDTO struct {
	EmployeeID domain.ID
}

type OrganizationResponsibleRepository interface {
	GetList(ctx context.Context, dto GetOrganizationResponsiblesListDTO) ([]domain.OrganizationResponsible, error)
	Get(ctx context.Context, dto GetOrganizationResponsibleDTO) (*domain.OrganizationResponsible, error)
}

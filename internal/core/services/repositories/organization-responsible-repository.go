package repositories

import (
	"context"
	"github.com/FFH255/avito-tender-management-golang/internal/core/domain"
)

type GetOneOrganizationResponsibleDTO struct {
	Username string
}

type OrganizationResponsibleRepository interface {
	GetOne(ctx context.Context, dto GetOneOrganizationResponsibleDTO) (domain.OrganizationResponsible, error)
}

package organization_responsible_repository

import (
	"tms/src/core/services/repositories"
	"tms/src/pkg/pg"
)

type OrganizationResponsibleRepository struct {
	client pg.Client
}

func New(client pg.Client) repositories.OrganizationResponsibleRepository {
	return OrganizationResponsibleRepository{client: client}
}

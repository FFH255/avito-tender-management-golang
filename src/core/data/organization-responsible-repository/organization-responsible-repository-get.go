package organization_responsible_repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

func (r OrganizationResponsibleRepository) Get(ctx context.Context, dto repositories.GetOrganizationResponsibleDTO) (*domain.OrganizationResponsible, error) {
	query := `SELECT id, organization_id, user_id FROM organization_responsible WHERE user_id = $1`

	args := []interface{}{dto.EmployeeID}

	row := r.client.QueryRow(ctx, query, args...)

	var orgResponsible domain.OrganizationResponsible

	err := row.Scan(&orgResponsible.ID, &orgResponsible.OrganizationID, &orgResponsible.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(domain.ErrNotFound, "orgResponsible not found")
		}
		return nil, err
	}

	return &orgResponsible, nil
}

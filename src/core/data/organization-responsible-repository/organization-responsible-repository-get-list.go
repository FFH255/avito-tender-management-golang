package organization_responsible_repository

import (
	"context"
	"fmt"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

func (r OrganizationResponsibleRepository) GetList(ctx context.Context, dto repositories.GetOrganizationResponsiblesListDTO) ([]domain.OrganizationResponsible, error) {
	query := `SELECT id, organization_id, user_id FROM organization_responsible WHERE 1=1`
	args := make([]interface{}, 0)
	i := 1

	if dto.OrganizationID != nil {
		args = append(args, *dto.OrganizationID)
		query += fmt.Sprintf(" AND organization_id = $%d", i)
		i++
	}

	if dto.EmployeeID != nil {
		args = append(args, *dto.EmployeeID)
		query += fmt.Sprintf(" AND user_id = $%d", i)
		i++
	}

	rows, err := r.client.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgResponsibles []domain.OrganizationResponsible

	for rows.Next() {
		var orgResponsible domain.OrganizationResponsible
		err := rows.Scan(&orgResponsible.ID, &orgResponsible.OrganizationID, &orgResponsible.UserID)
		if err != nil {
			return nil, err
		}
		orgResponsibles = append(orgResponsibles, orgResponsible)
	}

	return orgResponsibles, nil
}

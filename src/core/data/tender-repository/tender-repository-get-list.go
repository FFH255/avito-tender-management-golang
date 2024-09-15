package tender_repository

import (
	"context"
	"fmt"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

func (r TenderRepository) GetList(ctx context.Context, dto repositories.GetTendersListDTO) ([]domain.Tender, error) {
	query := `SELECT id, name, description, service_type, status, organization_id, version FROM tender WHERE 1=1`
	args := make([]interface{}, 0)
	i := 1

	if dto.OrganizationID != nil {
		query += fmt.Sprintf(` AND organization_id = $%d`, i)
		args = append(args, dto.OrganizationID)
		i++
	}

	if dto.ServiceType != nil {
		query += fmt.Sprintf(` AND service_type = $%d`, i)
		args = append(args, dto.ServiceType)
		i++
	}

	if dto.Status != nil {
		query += fmt.Sprintf(` AND status = $%d`, i)
		args = append(args, dto.Status)
		i++
	}

	if dto.Limit != nil {
		query += fmt.Sprintf(` LIMIT $%d`, i)
		args = append(args, dto.Limit)
		i++
	}

	if dto.Offset != nil {
		query += fmt.Sprintf(` OFFSET $%d`, i)
		args = append(args, dto.Offset)
		i++
	}

	rows, err := r.client.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenders := make([]domain.Tender, 0)

	for rows.Next() {
		var tender domain.Tender

		err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version)

		if err != nil {
			return nil, err
		}

		query = `SELECT name, description, service_type, version, created_at FROM tender_snapshot WHERE tender_id = $1`

		snapshotRows, err := r.client.Query(ctx, query, tender.ID)

		if err != nil {
			return nil, err
		}
		defer snapshotRows.Close()

		snapshots := make([]domain.TenderSnapshot, 0)

		for snapshotRows.Next() {
			var snapshot domain.TenderSnapshot

			err := snapshotRows.Scan(&snapshot.Name, &snapshot.Description, &snapshot.ServiceType, &snapshot.Version, &snapshot.CreatedAt)

			if err != nil {
				return nil, err
			}
			snapshots = append(snapshots, snapshot)
		}
		tender.Snapshots = snapshots
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

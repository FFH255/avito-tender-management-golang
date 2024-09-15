package tender_repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

func (r TenderRepository) Get(ctx context.Context, dto repositories.GetTenderDTO) (*domain.Tender, error) {
	query := `SELECT id, name, description, service_type, status, organization_id, version FROM tender WHERE id = $1`
	args := []interface{}{dto.ID}
	i := 2

	if dto.OrganizationID != nil {
		query += fmt.Sprintf(` AND organization_id = $%d`, i)
		args = append(args, dto.OrganizationID)
		i++
	}

	row := r.client.QueryRow(ctx, query, args...)

	var tender domain.Tender

	err := row.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(domain.ErrNotFound, "tender not found")
		}
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

	return &tender, nil
}

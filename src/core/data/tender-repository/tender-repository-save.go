package tender_repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"tms/src/core/domain"
)

func (r TenderRepository) Save(ctx context.Context, tender domain.Tender) error {

	const (
		deleteTenderQuery = `DELETE FROM tender WHERE id=$1`

		deleteTenderSnapshotsQuery = `DELETE FROM tender_snapshot WHERE tender_id=$1`

		createTenderQuery = `INSERT INTO tender(id, name, description, service_type, status, organization_id, version, created_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

		createTenderSnapshotQuery = `INSERT INTO tender_snapshot(id, tender_id, name, description, service_type, version, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7);`
	)

	tx, err := r.client.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, deleteTenderQuery, tender.ID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteTenderSnapshotsQuery, tender.ID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, createTenderQuery, tender.ID, tender.Name, tender.Description, tender.ServiceType,
		tender.Status, tender.OrganizationID, tender.Version, tender.CreatedAt)

	if err != nil {
		return err
	}

	for _, snapshot := range tender.Snapshots {
		_, err = tx.Exec(ctx, createTenderSnapshotQuery, snapshot.ID, tender.ID, snapshot.Name, snapshot.Description,
			snapshot.ServiceType, snapshot.Version, snapshot.CreatedAt)

		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

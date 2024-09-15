package bid_repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"tms/src/core/domain"
)

func (r BidRepository) Save(ctx context.Context, bid domain.Bid) error {
	tx, err := r.client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const (
		deleteBid = `DELETE FROM bid WHERE id = $1`

		deleteSnapshots = `DELETE FROM bid_snapshot WHERE bid_id = $1`

		insertBid = `INSERT INTO bid(ID, NAME, DESCRIPTION, STATUS, TENDER_ID, AUTHOR_TYPE, AUTHOR_ID, VERSION, CREATED_AT) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

		insertSnapshot = `INSERT INTO bid_snapshot(id, bid_id, name, description, version) VALUES ($1, $2, $3, $4, $5)`
	)

	_, err = tx.Exec(ctx, deleteBid, bid.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteSnapshots, bid.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, insertBid, bid.ID, bid.Name, bid.Description, bid.Status, bid.TenderID, bid.AuthorType, bid.AuthorID, bid.Version, bid.CreatedAt)
	if err != nil {
		return err
	}

	for _, s := range bid.Snapshots {
		_, err = tx.Exec(ctx, insertSnapshot, s.ID, bid.ID, s.Name, s.Description, s.Version)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

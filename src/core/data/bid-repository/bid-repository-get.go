package bid_repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

func (r BidRepository) Get(ctx context.Context, dto repositories.GetBidDTO) (*domain.Bid, error) {
	queryBid := `SELECT id, name, description, status, tender_id, author_type, author_id, version, created_at FROM bid WHERE id = $1`

	querySnapshots := `SELECT id, name, description, version FROM bid_snapshot WHERE bid_id = $1`

	args := []interface{}{dto.ID}
	i := 2

	if dto.AuthorID != nil {
		args = append(args, dto.AuthorID)
		queryBid += fmt.Sprintf(" AND author_id = $%d", i)
		i++
	}
	row := r.client.QueryRow(ctx, queryBid, args...)

	var bid domain.Bid

	err := row.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(domain.ErrNotFound, "bid not found")
		}
		return nil, err
	}

	rows, err := r.client.Query(ctx, querySnapshots, bid.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s domain.BidSnapshot
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.Version); err != nil {
			return nil, err
		}
		bid.Snapshots = append(bid.Snapshots, s)
	}

	return &bid, nil
}

package bid_repository

import (
	"context"
	"fmt"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

func (r BidRepository) GetList(ctx context.Context, dto repositories.GetBidListDTO) ([]domain.Bid, error) {
	var (
		queryBids = `SELECT id, name, description, status, tender_id, author_type, author_id, version, created_at FROM bid WHERE 1=1`

		querySnapshots = `SELECT id, name, description, version FROM bid_snapshot WHERE bid_id = $1`
	)

	args := make([]interface{}, 0)
	i := 1

	if dto.ID != nil {
		args = append(args, string(*dto.ID))
		queryBids += fmt.Sprintf(" AND id = $%d", i)
		i++
	}

	if dto.Status != nil {
		args = append(args, *dto.Status)
		queryBids += fmt.Sprintf(" AND status = $%d", i)
		i++
	}

	if dto.AuthorType != nil {
		args = append(args, *dto.AuthorType)
		queryBids += fmt.Sprintf(" AND author_type = $%d", i)
		i++
	}

	if dto.AuthorID != nil {
		args = append(args, *dto.AuthorID)
		queryBids += fmt.Sprintf(" AND author_id = $%d", i)
		i++
	}

	if dto.TenderID != nil {
		args = append(args, *dto.TenderID)
		queryBids += fmt.Sprintf(" AND tender_id = $%d", i)
		i++
	}

	if dto.Limit != nil {
		args = append(args, *dto.Limit)
		queryBids += fmt.Sprintf(" LIMIT $%d", i)
		i++
	}

	if dto.Offset != nil {
		args = append(args, *dto.Offset)
		queryBids += fmt.Sprintf(" OFFSET $%d", i)
		i++
	}

	rows, err := r.client.Query(ctx, queryBids, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []domain.Bid

	for rows.Next() {
		var bid domain.Bid
		err := rows.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt)
		if err != nil {
			return nil, err
		}

		snapshotRows, err := r.client.Query(ctx, querySnapshots, bid.ID)
		if err != nil {
			return nil, err
		}
		defer snapshotRows.Close()

		for snapshotRows.Next() {
			var s domain.BidSnapshot
			if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.Version); err != nil {
				return nil, err
			}
			bid.Snapshots = append(bid.Snapshots, s)
		}
		bids = append(bids, bid)
	}

	return bids, nil
}

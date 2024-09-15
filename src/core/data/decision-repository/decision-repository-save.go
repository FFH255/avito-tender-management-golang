package decision_repository

import (
	"context"
	"tms/src/core/domain"
)

func (r DecisionRepository) Save(ctx context.Context, decision domain.Decision) error {
	query := `INSERT INTO decision(id, author_id, bid_id, tender_id, status) VALUES($1, $2, $3, $4, $5)`
	_, err := r.client.Exec(ctx, query, decision.ID, decision.AuthorID, decision.BidID, decision.TenderID, decision.Status)
	return err
}

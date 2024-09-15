package decision_repository

import (
	"context"
	"fmt"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

func (r DecisionRepository) GetList(ctx context.Context, dto repositories.GetDecisionListDTO) ([]domain.Decision, error) {
	query := `SELECT id, author_id, bid_id, tender_id, status FROM decision WHERE 1=1`
	args := make([]interface{}, 0)
	i := 1

	if dto.TenderID != nil {
		args = append(args, *dto.TenderID)
		query += fmt.Sprintf(" AND tender_id = $%d", i)
		i++
	}

	rows, err := r.client.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decisions []domain.Decision

	for rows.Next() {
		var decision domain.Decision
		if err := rows.Scan(&decision.ID, &decision.AuthorID, &decision.BidID, &decision.TenderID, &decision.Status); err != nil {
			return nil, err
		}
		decisions = append(decisions, decision)
	}

	return decisions, nil
}

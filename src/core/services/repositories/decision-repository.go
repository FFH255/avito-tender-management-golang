package repositories

import (
	"context"
	"tms/src/core/domain"
)

type GetDecisionListDTO struct {
	TenderID *domain.ID
}

type DecisionRepository interface {
	GetList(ctx context.Context, dto GetDecisionListDTO) ([]domain.Decision, error)
	Save(ctx context.Context, decision domain.Decision) error
}

package repositories

import (
	"context"
	"github.com/FFH255/avito-tender-management-golang/internal/core/domain"
)

type GetOneTenderDTO struct {
	ID string
}

type TenderRepository interface {
	GetOne(ctx context.Context, dto GetOneTenderDTO) (domain.Tender, error)
	Save(ctx context.Context, tender domain.Tender) error
}

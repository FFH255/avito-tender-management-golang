package repositories

import (
	"context"
	"tms/src/core/domain"
)

type GetBidListDTO struct {
	ID         *domain.ID
	Status     *domain.BidStatus
	TenderID   *domain.ID
	AuthorType *domain.BidAuthorType
	AuthorID   *domain.ID
	Limit      *Limit
	Offset     *Offset
}

type GetBidDTO struct {
	ID       domain.ID
	AuthorID *domain.ID
}

type BidRepository interface {
	GetList(ctx context.Context, dto GetBidListDTO) ([]domain.Bid, error)
	Get(ctx context.Context, dto GetBidDTO) (*domain.Bid, error)
	Save(ctx context.Context, bid domain.Bid) error
}

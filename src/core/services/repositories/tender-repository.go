package repositories

import (
	"context"
	"tms/src/core/domain"
)

type GetTendersListDTO struct {
	OrganizationID *domain.ID
	Status         *domain.TenderStatus
	ServiceType    *domain.TenderServiceType
	Offset         *Offset
	Limit          *Limit
}

type GetTenderDTO struct {
	ID             domain.ID
	OrganizationID *domain.ID
}

type TenderRepository interface {
	GetList(ctx context.Context, dto GetTendersListDTO) ([]domain.Tender, error)
	Get(ctx context.Context, dto GetTenderDTO) (*domain.Tender, error)
	Save(ctx context.Context, tender domain.Tender) error
}

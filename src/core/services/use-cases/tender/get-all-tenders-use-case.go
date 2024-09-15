package use_cases

import (
	"context"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type GetAllTendersDTO struct {
	Limit       *int    `json:"limit"`
	Offset      *int    `json:"offset"`
	ServiceType *string `json:"service_type"`
}

type GetAllTendersUseCase struct {
	tenderRepository repositories.TenderRepository
}

func (uc GetAllTendersUseCase) Execute(dto GetAllTendersDTO) ([]domain.Tender, error) {
	limit := repositories.NewLimit(dto.Limit)
	offset := repositories.NewOffset(dto.Offset)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var serviceType *domain.TenderServiceType

	if dto.ServiceType != nil {
		st, err := domain.NewTenderServiceType(*dto.ServiceType)
		if err != nil {
			return nil, err
		}
		serviceType = &st
	}

	status := domain.TenderPublishedStatus
	tenders, err := uc.tenderRepository.GetList(ctx, repositories.GetTendersListDTO{
		ServiceType: serviceType,
		Offset:      &offset,
		Limit:       &limit,
		Status:      &status,
	})

	if err != nil {
		return nil, err
	}

	return tenders, nil
}

func NewGetAllTendersUseCase(tenderRepository repositories.TenderRepository) GetAllTendersUseCase {
	return GetAllTendersUseCase{
		tenderRepository: tenderRepository,
	}
}

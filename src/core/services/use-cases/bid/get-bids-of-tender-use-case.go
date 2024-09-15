package use_cases

import (
	"context"
	"github.com/pkg/errors"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type GetBidsOfTenderUseCase struct {
	employeeRepository       repositories.EmployeeRepository
	orgResponsibleRepository repositories.OrganizationResponsibleRepository
	tenderRepository         repositories.TenderRepository
	bidRepository            repositories.BidRepository
}

func NewGetBidsOfTenderUseCase(
	employeeRepository repositories.EmployeeRepository,
	tenderRepository repositories.TenderRepository,
	bidRepository repositories.BidRepository,
	orgResponsibleRepository repositories.OrganizationResponsibleRepository,
) GetBidsOfTenderUseCase {
	return GetBidsOfTenderUseCase{
		employeeRepository:       employeeRepository,
		tenderRepository:         tenderRepository,
		bidRepository:            bidRepository,
		orgResponsibleRepository: orgResponsibleRepository,
	}
}

type GetBidsOfTenderDTO struct {
	TenderID string
	Username string
	Limit    *int
	Offset   *int
}

func (uc GetBidsOfTenderUseCase) Execute(dto GetBidsOfTenderDTO) ([]domain.Bid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Проверка существования Employee
	employee, err := uc.employeeRepository.Get(ctx, repositories.GetEmployeeDTO{
		Username: &dto.Username,
	})
	if err != nil {
		return nil, err
	}

	// Проверка существования OrgResponsible
	orgResp, err := uc.orgResponsibleRepository.Get(ctx, repositories.GetOrganizationResponsibleDTO{
		EmployeeID: employee.ID,
	})
	if err != nil {
		return nil, err
	}

	// Проверка существования Tender
	tenderID := domain.ID(dto.TenderID)
	tender, err := uc.tenderRepository.Get(ctx, repositories.GetTenderDTO{
		ID: tenderID,
	})
	if err != nil {
		return nil, err
	}

	// Проверка прав OrgReps
	if orgResp.OrganizationID != tender.OrganizationID {
		return nil, errors.Wrap(domain.ErrNoPermission, "organization does not belong to tender")
	}

	// Получение списка Bid
	limit := repositories.NewLimit(dto.Limit)
	offset := repositories.NewOffset(dto.Offset)
	bids, err := uc.bidRepository.GetList(ctx, repositories.GetBidListDTO{
		TenderID: &tender.ID,
		Limit:    &limit,
		Offset:   &offset,
	})
	if err != nil {
		return nil, err
	}

	return bids, nil
}

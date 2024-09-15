package use_cases

import (
	"context"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type ChangeTenderStatusUseCase struct {
	employeeRepository       repositories.EmployeeRepository
	orgResponsibleRepository repositories.OrganizationResponsibleRepository
	tenderRepository         repositories.TenderRepository
}

type ChangeTenderStatusDTO struct {
	TenderID string
	Status   string
	Username string
}

func (uc *ChangeTenderStatusUseCase) Execute(dto ChangeTenderStatusDTO) (*domain.Tender, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	employee, err := uc.employeeRepository.Get(ctx, repositories.GetEmployeeDTO{
		Username: &dto.Username,
	})

	if err != nil {
		return nil, err
	}

	orgResp, err := uc.orgResponsibleRepository.Get(ctx, repositories.GetOrganizationResponsibleDTO{
		EmployeeID: employee.ID,
	})

	if err != nil {
		return nil, err
	}

	tenderID := domain.ID(dto.TenderID)

	tender, err := uc.tenderRepository.Get(ctx, repositories.GetTenderDTO{
		ID: tenderID,
	})

	if err != nil {
		return nil, err
	}

	if err = tender.ChangeStatus(*orgResp, dto.Status); err != nil {
		return nil, err
	}

	if err = uc.tenderRepository.Save(ctx, *tender); err != nil {
		return nil, err
	}

	return tender, nil
}

func NewChangeTenderStatusUseCase(
	employeeRepository repositories.EmployeeRepository,
	orgResponsibleRepository repositories.OrganizationResponsibleRepository,
	tenderRepository repositories.TenderRepository,
) ChangeTenderStatusUseCase {
	return ChangeTenderStatusUseCase{
		employeeRepository:       employeeRepository,
		orgResponsibleRepository: orgResponsibleRepository,
		tenderRepository:         tenderRepository,
	}
}

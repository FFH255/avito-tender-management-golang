package use_cases

import (
	"context"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type GetTenderStatusUseCase struct {
	employeeRepository       repositories.EmployeeRepository
	orgResponsibleRepository repositories.OrganizationResponsibleRepository
	tenderRepository         repositories.TenderRepository
}

type GetTenderStatusDTO struct {
	TenderID string
	Username string
}

func (uc GetTenderStatusUseCase) Execute(dto GetTenderStatusDTO) (domain.TenderStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	employee, err := uc.employeeRepository.Get(ctx, repositories.GetEmployeeDTO{
		Username: &dto.Username,
	})
	if err != nil {
		return "", err
	}

	orgResponsible, err := uc.orgResponsibleRepository.Get(ctx, repositories.GetOrganizationResponsibleDTO{
		EmployeeID: employee.ID,
	})
	if err != nil {
		return "", err
	}

	tenderID := domain.ID(dto.TenderID)

	tender, err := uc.tenderRepository.Get(ctx, repositories.GetTenderDTO{
		ID:             tenderID,
		OrganizationID: &orgResponsible.OrganizationID,
	})
	if err != nil {
		return "", err
	}

	return tender.Status, nil
}

func NewGetTenderStatusUseCase(
	employeeRepository repositories.EmployeeRepository,
	orgResponsibleRepository repositories.OrganizationResponsibleRepository,
	tenderRepository repositories.TenderRepository,
) GetTenderStatusUseCase {
	return GetTenderStatusUseCase{
		employeeRepository:       employeeRepository,
		orgResponsibleRepository: orgResponsibleRepository,
		tenderRepository:         tenderRepository,
	}
}

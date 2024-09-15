package use_cases

import (
	"context"
	"time"
	"tms/src/core/domain"
	"tms/src/core/services/repositories"
)

type RollBackTenderUseCase struct {
	employeeRepository       repositories.EmployeeRepository
	orgResponsibleRepository repositories.OrganizationResponsibleRepository
	tenderRepository         repositories.TenderRepository
}

type RollBackTenderUseCaseDTO struct {
	TenderID string
	Version  int
	Username string
}

func (uc RollBackTenderUseCase) Execute(dto RollBackTenderUseCaseDTO) (*domain.Tender, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	employee, err := uc.employeeRepository.Get(ctx, repositories.GetEmployeeDTO{
		Username: &dto.Username,
	})
	if err != nil {
		return nil, err
	}

	orgResponsible, err := uc.orgResponsibleRepository.Get(ctx, repositories.GetOrganizationResponsibleDTO{
		EmployeeID: employee.ID,
	})
	if err != nil {
		return nil, err
	}

	tenderID := domain.ID(dto.TenderID)
	tender, err := uc.tenderRepository.Get(ctx, repositories.GetTenderDTO{
		OrganizationID: &orgResponsible.OrganizationID,
		ID:             tenderID,
	})
	if err != nil {
		return nil, err
	}

	if err := tender.Rollback(*orgResponsible, dto.Version); err != nil {
		return nil, err
	}

	if err = uc.tenderRepository.Save(ctx, *tender); err != nil {
		return nil, err
	}

	return tender, nil
}

func NewRollBackTenderUseCase(
	employeeRepository repositories.EmployeeRepository,
	orgResponsibleRepository repositories.OrganizationResponsibleRepository,
	tenderRepository repositories.TenderRepository,
) RollBackTenderUseCase {
	return RollBackTenderUseCase{
		employeeRepository:       employeeRepository,
		orgResponsibleRepository: orgResponsibleRepository,
		tenderRepository:         tenderRepository,
	}
}
